package encoding

import (
	"ehealth-migration/application"
	"ehealth-migration/pkg/config"
	"ehealth-migration/pkg/database"
	"ehealth-migration/pkg/generated"
	"ehealth-migration/pkg/pools"
	"ehealth-migration/pkg/reflections"
	"ehealth-migration/pkg/structtag"
	"ehealth-migration/pkg/trait"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"time"
	//"strings"
)

type Encoding struct {
	app   *application.Application
	Model interface{}
}

func (e *Encoding) GetModel() interface{} {
	return e.Model
}

func (e *Encoding) SetModel(fields []config.Field) {
	var generatedStruct generated.GeneratedStruct
	var filedsStruct []generated.Field

	for _, filed := range fields {
		filedsStruct = append(filedsStruct, generated.Field{
			Name: filed.Name,
			Type: filed.Type,
			Tag:  filed.Tag,
		})
	}

	generatedStruct.Fields = filedsStruct
	e.Model = generatedStruct.ToStruct()
}

func Run(app *application.Application) error {
	var err error
	log.Debug().Msgf("Run encoding")

	encoding := Encoding{app: app}

	flags := encoding.app.Flags
	config := encoding.app.Config

	if !config.IsDatabase(flags.Name) {
		return errors.New(fmt.Sprintf("Can not find Database %s", flags.Name))
	}

	currentDatabase := config.GetDatabase(flags.Name)
	if !currentDatabase.IsTable(flags.Tname) {
		return errors.New(fmt.Sprintf("Can not find Table %s", flags.Tname))
	}

	currentTable := currentDatabase.GetTable(flags.Tname)
	if err := currentTable.Validate(); err != nil {
		return err
	}

	if filedErrors := currentTable.CheckFields(); len(filedErrors) != 0 {
		log.Error().Errs("Init CheckFields", filedErrors).Msg("encoding.CheckFields")
		return errors.New("CheckFields errors")
	}

	// generated dynamic struct
	encoding.SetModel(currentTable.Fields)

	// stup
	// indexes
	err = encoding.app.DbClient.CreateTableIndexes(currentTable.Name)
	if err != nil {
		log.Error().Err(err).Msg("CreateTableIndexes")
		return err
	}

	err = encoding.SaveAllIndexes(currentTable)
	if err != nil {
		return err
	}

	err = encoding.DropAllIndexes(currentTable)
	if err != nil {
		return err
	}

	// truncate
	truncates := currentDatabase.GetTruncates()
	for _, truncate := range truncates {
		errTruncate := encoding.app.DbClient.Truncate(truncate.Name)
		if errTruncate != nil {
			log.Error().Err(errTruncate).Msg("Truncate error")
			return errTruncate
		}
	}

	// created Table
	err = encoding.app.DbClient.CreateStatusEncodingTable(currentTable.Name)
	if err != nil {
		log.Error().Err(err).Msg("CreateStatusEncodingTable")
		return err
	}

	err = encoding.app.DbClient.SetStatusEncodingDefaultValues(currentTable.Name)
	if err != nil {
		log.Error().Err(err).Msg("SetStatusEncodingDefaultValues")
		return err
	}

	err = encoding.app.DbClient.CreateFailedTables(currentTable.Name)
	if err != nil {
		log.Error().Err(err).Msg("CreateFailedTables")
		return err
	}

	encoding.app.DbClient.Migrate(currentTable.TmpName, encoding.GetModel())

	// Pools Run
	status, err := encoding.PoolsTmp(currentTable)
	if err != nil {
		return err
	}

	// processing Done
	if status {
		log.Info().Msg("Start processing rename table and created index")

		newNameTable := currentTable.Name + "_" + trait.Int64ToString(trait.MakeTimestamp())
		log.Info().Msgf("Rename table: %s -> %s", currentTable.Name, newNameTable)
		err = encoding.app.DbClient.RenameTable(currentTable.Name, newNameTable)
		if err != nil {
			return err
		}

		log.Info().Msgf("Rename table: %s -> %s", currentTable.TmpName, currentTable.Name)
		err = encoding.app.DbClient.RenameTable(currentTable.TmpName, currentTable.Name)
		if err != nil {
			return err
		}

		log.Info().Msg("CreatedAllIndexes")
		// created all indexes
		err = encoding.CreatedAllIndexes(currentTable)
		if err != nil {
			return err
		}

		log.Info().Msg("Done")
	}

	return nil
}

// ==================== Set Table Tmp ============
func (e *Encoding) PoolsTmp(currentTable config.Table) (bool, error) {

	statusDone := false
	iteration := 0

	status, err := e.app.DbClient.GetOffset(currentTable.Name)
	if err != nil {
		log.Error().Err(err).Msg("DbClient.GetOffset error")
		return statusDone, err
	}
	lasId, err := e.app.DbClient.GetTmpLasInsert(currentTable.TmpName)
	if err != nil {
		log.Error().Err(err).Msg("DbClient.GetTmpLasInsert error")
		return statusDone, err
	}

	for {
		var itemsLen int
		tasks := make([]*pools.Task, 0)
		models := generated.MakeSlices(e.GetModel())

		log.Info().Msgf("Run iterations (%d) - (%d / %d) \n", iteration, status.Result, e.app.Flags.Limit)

		items, err := e.app.DbClient.GetList(currentTable.Name, models, lasId, e.app.Flags.Limit)
		if err != nil {
			log.Error().Err(err).Msg("GetList error")
			break
		}

		//// TO DO testing
		//if iteration == 1 {
		//	break
		//}

		tasks, itemsLen, lasId = e.getTasks(items, currentTable)
		if itemsLen == 0 {
			statusDone = true
			break
		}

		if err := e.initPools(tasks); err != nil {
			log.Fatal().Msgf("Error Pools tasks `%s`: %s \n", currentTable.Name, err)
		}

		status.Result += itemsLen

		iteration++
	}

	return statusDone, nil
}

// ================= Init Pools ============
func (e *Encoding) getTasks(result interface{}, table config.Table) ([]*pools.Task, int, string) {
	var lenTasks int
	var lastId string
	tasks := make([]*pools.Task, 0)

	// get values of interface
	if reflections.IsSlice(result) || reflections.IsPointer(result) {
		collection := reflections.GetReflectValue(result)

		lenTasks = collection.Len()
		for i := 0; i < lenTasks; i++ {

			r := collection.Index(i)
			if reflections.IsStruct(r) {

				obj := r.Interface()

				// get tag generated
				obj, err := generatedData(obj, table.Fields, table.Lang)
				if err != nil {
					log.Fatal().Err(err).Msg("generatedData error")
				}

				value, err := reflections.GetField(obj, "ID")
				if err != nil {
					log.Fatal().Err(err).Msg("reflections.GetField error")
				}

				currentID := value.(uuid.UUID).String()

				// TO DO New Tasks
				tasks = append(tasks, pools.NewTask(func(t *pools.Task) error {
					e := t.App.(*Encoding)
					object := t.Message.Object
					table := t.Message.Table.(config.Table)

					log.Debug().Msgf("Run task message[%+s]", t.Message.ID)

					err = e.app.DbClient.InsertTmp(table.TmpName, object)
					if err != nil {
						log.Error().Err(err).Msg("InsertTmp error")
						uuid, err := uuid.FromString(t.Message.ID)
						if err != nil {
							log.Error().Err(err).Msg("uuid.FromString")
							return err
						}

						err = e.app.DbClient.SetFailed(table.Name, database.IFailed{Id: uuid, CreatedAt: time.Now()})
						if err != nil {
							log.Error().Err(err).Msg("SetFailed")
						}
						return err
					}

					err = e.app.DbClient.UpdateOffset(table.Name)
					if err != nil {
						log.Fatal().Msgf("Fatal UpdateOffset `%s`: %s \n", "declaration_requests", err)
						return err
					}

					return nil
				}, &pools.Message{Object: obj, Table: table, ID: currentID}, e))

				if i == lenTasks-1 {
					lastId = currentID
				}

			}
		}
	}

	return tasks, lenTasks, lastId
}

func (e *Encoding) initPools(tasks []*pools.Task) error {

	p, err := pools.NewPool(
		tasks,
		e.app.Flags.Concurrency,
		100)
	if err != nil {
		return err
	}
	p.Run()

	for _, task := range p.Tasks {
		if task.Err != nil {
			log.Error().Err(task.Err).Msg("Task error msg.")
			continue
		}
	}

	return nil
}

// encoding fields
var (
	TAG_NAME = "generated"
)

func generatedData(object interface{}, fields []config.Field, lang string) (interface{}, error) {
	for _, field := range fields {
		value, err := reflections.GetField(object, field.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("reflections.GetField error")
		}

		tag, err := reflections.GetFieldTag(object, field.Name, TAG_NAME)
		if err != nil {
			log.Fatal().Err(err).Msg("reflections.GetFieldTag error")
		}

		if tag != "" {
			currentTag := structtag.Parse(tag)

			if currentTag.Key != "" {
				valueData := generated.GetGenerated(currentTag, value, lang)
				if valueData == nil {
					return nil, errors.New("generated.GetGenerated not find tag: " + currentTag.Nmae)
				}

				err := reflections.SetField(object, field.Name, valueData)
				if err != nil {
					log.Fatal().Err(err).Msg("reflections.SetField error")
				}
			}
		}
	}

	return object, nil
}

// Indexes
func (e *Encoding) DropAllIndexes(currentTable config.Table) error {
	indexes, err := e.app.DbClient.GetAllIndexes(currentTable.Name)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		err := e.app.DbClient.DropIndex(index.Indexname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoding) CreatedAllIndexes(currentTable config.Table) error {
	indexes, err := e.app.DbClient.GetAllIndexes(currentTable.Name)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		err := e.app.DbClient.CreatedIndex(index.Indexdef)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoding) SaveAllIndexes(currentTable config.Table) error {
	indexes, err := e.app.DbClient.GetCurrentTableIndexs(currentTable.Name, currentTable.PkeyName)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		cIndex, err := e.app.DbClient.GetIndex(currentTable.Name, index.Indexname)
		if err != nil {
			return err
		}

		if cIndex.Id == 0 {
			err := e.app.DbClient.SetIndex(currentTable.Name, &index)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
