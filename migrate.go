package uadmin

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/uadmin/uadmin/core"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Migration struct {
	gorm.Model
	MigrationName string `gorm:"index:migration_migration_name,unique"`
	AppliedAt     time.Time
}

type MigrateCommand struct {
}

func (c MigrateCommand) Proceed(subaction string, args []string) error {
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.addAction("create", &CreateMigration{})
	commandRegistry.addAction("up", &UpMigration{})
	commandRegistry.addAction("down", &DownMigration{})
	commandRegistry.addAction("determine-conflicts", &DetermineConflictsMigration{})
	isCorrectActionPassed = commandRegistry.isRegisteredCommand(subaction)
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	return commandRegistry.runAction(subaction, "", args)
}

func (c MigrateCommand) GetHelpText() string {
	return "Migrate your database"
}

func prepareMigrationName(message string) string {
	now := time.Now()
	sec := now.Unix()
	message = core.AsciiRegex.ReplaceAllLiteralString(message, "")
	if len(message) > 30 {
		message = message[:30]
	}
	message = strings.Replace(strings.ToLower(message), " ", "_", -1)
	message = strings.Replace(strings.ToLower(message), ".", "_", -1)
	return fmt.Sprintf("%s_%d", message, sec)
}

type CreateMigrationOptions struct {
	Message   string `short:"m" description:"Describe what is this migration for"`
	Blueprint string `short:"b" description:"Blueprint you'd like to create migration for"`
	MergeMode bool   `long:"merge" description:"Merge conflicted migrations"`
}

type CreateMigration struct {
}

func (command CreateMigration) Proceed(subaction string, args []string) error {
	var opts = &CreateMigrationOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if len(args) == 0 {
		var help string = `
Please provide flags -b and -m which are blueprint and description of the migration respectively
or --merge if you want to merge conflicted migrations
`
		fmt.Printf(help)
		return nil
	}

	if !opts.MergeMode && (opts.Blueprint == "" || opts.Message == "") {
		var help string = `
Please provide flags -b and -m which are blueprint and description of the migration respectively or --merge
if you want to merge conflicted migrations 
`
		fmt.Printf(help)
		return nil
	}

	if err != nil {
		return err
	}
	const concreteMigrationTpl = `package migrations

type {{.MigrationName}} struct {
}

func (m {{.MigrationName}}) GetName() string {
    return "{{.BlueprintName}}.{{.ConcreteMigrationId}}"
}

func (m {{.MigrationName}}) GetId() int64 {
    return {{.ConcreteMigrationId}}
}

func (m {{.MigrationName}}) Up(uadminDatabase *core.UadminDatabase) error {
    return nil
}

func (m {{.MigrationName}}) Down(uadminDatabase *core.UadminDatabase) error {
    return nil
}

func (m {{.MigrationName}}) Deps() []string {
{{if .Dependencies}}    return []string{{.Dependencies}}{{else}}    return make([]string, 0){{end}}
}
`
	const initializeMigrationRegistryTpl = `
    BMigrationRegistry.AddMigration({{.MigrationName}}{})`
	const migrationRegistryCreationTpl = `package migrations

import (
	"github.com/uadmin/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
    BMigrationRegistry = core.NewMigrationRegistry()
    // placeholder to insert next migration
}
`
	if opts.MergeMode {
		for true {
			for traverseMigrationResult := range appInstance.BlueprintRegistry.TraverseMigrations() {
				if traverseMigrationResult.Error != nil && strings.Contains(traverseMigrationResult.Error.Error(), "Found two or more migrations with no children from the same blueprint") {
					r := regexp.MustCompile(`Set\{(.*?)\}`)
					conflictedMigrations := r.FindStringSubmatch(traverseMigrationResult.Error.Error())[1]
					listOfConflictedMigrations := strings.Split(conflictedMigrations, ",")
					dependenciesString := make([]string, len(listOfConflictedMigrations))
					for i, conflictedMigration := range listOfConflictedMigrations {
						listOfConflictedMigrations[i] = strings.TrimSpace(conflictedMigration)
						dependenciesString[i] = fmt.Sprintf(`"%s"`, listOfConflictedMigrations[i])
					}
					migrationName := prepareMigrationName(strings.Join(listOfConflictedMigrations, "_"))
					blueprintName := core.GetBluePrintNameFromMigrationName(listOfConflictedMigrations[0])
					dirPath := "blueprint/" + strings.ToLower(blueprintName) + "/migrations"
					pathToBaseMigrationsFile := dirPath + "/migrations.go"
					pathToConcreteMigrationsFile := dirPath + "/" + migrationName + ".go"
					var concreteTplBuffer bytes.Buffer
					now := time.Now()
					sec := now.Unix()
					concreteTpl := template.Must(template.New("concretemigration").Parse(concreteMigrationTpl))
					concreteData := struct {
						MigrationName       string
						ConcreteMigrationId string
						Dependencies        string
						BlueprintName       string
					}{
						MigrationName:       migrationName,
						ConcreteMigrationId: strconv.Itoa(int(sec)),
						Dependencies:        "{" + strings.Join(dependenciesString, ",") + "}",
						BlueprintName:       blueprintName,
					}
					if err = concreteTpl.Execute(&concreteTplBuffer, concreteData); err != nil {
						panic(err)
					}
					err = ioutil.WriteFile(pathToConcreteMigrationsFile, concreteTplBuffer.Bytes(), 0755)
					if err != nil {
						panic(err)
					}
					integrateMigrationIntoRegistryTpl := template.Must(template.New("integratemigrationintoregistry").Parse(initializeMigrationRegistryTpl))
					integrateMigrationIntoRegistryData := struct {
						MigrationName string
					}{
						MigrationName: migrationName,
					}
					var integrateMigrationIntoRegistryTplBuffer bytes.Buffer
					if err = integrateMigrationIntoRegistryTpl.Execute(&integrateMigrationIntoRegistryTplBuffer, integrateMigrationIntoRegistryData); err != nil {
						panic(err)
					}
					read, err := ioutil.ReadFile(pathToBaseMigrationsFile)
					if err != nil {
						panic(err)
					}
					newContents := strings.Replace(
						string(read),
						"// placeholder to insert next migration",
						integrateMigrationIntoRegistryTplBuffer.String()+"\n    // placeholder to insert next migration", -1)
					err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(newContents), 0755)
					if err != nil {
						panic(err)
					}
					fmt.Printf(
						"Created migration for blueprint %s with name %s\n",
						blueprintName,
						traverseMigrationResult.Error.Error(),
					)
				}
			}
			break
		}
	} else {
		bluePrintPath := "blueprint/" + strings.ToLower(opts.Blueprint)
		if _, err := os.Stat(bluePrintPath); os.IsNotExist(err) {
			panic(fmt.Sprintf("Blueprint %s doesn't exist", opts.Blueprint))
		}
		dirPath := "blueprint/" + strings.ToLower(opts.Blueprint) + "/migrations"
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err = os.Mkdir(dirPath, 0755)
			if err != nil {
				panic(err)
			}
		}
		pathToBaseMigrationsFile := dirPath + "/migrations.go"
		if _, err := os.Stat(pathToBaseMigrationsFile); os.IsNotExist(err) {
			err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(migrationRegistryCreationTpl), 0755)
			if err != nil {
				panic(err)
			}
		}
		migrationName := prepareMigrationName(opts.Message)
		pathToConcreteMigrationsFile := dirPath + "/" + migrationName + ".go"
		var lastMigrationId int
		err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			var migrationFileRegex = regexp.MustCompile(`.*?_(\d+)\.go`)
			match := migrationFileRegex.FindStringSubmatch(path)
			if len(match) > 0 {
				migrationId, _ := strconv.Atoi(match[1])
				if migrationId > lastMigrationId {
					lastMigrationId = migrationId
				}
			}
			return nil
		})
		var concreteTplBuffer bytes.Buffer
		now := time.Now()
		sec := now.Unix()
		concreteTpl := template.Must(template.New("concretemigration").Parse(concreteMigrationTpl))
		concreteData := struct {
			MigrationName       string
			ConcreteMigrationId string
			Dependencies        string
			BlueprintName       string
		}{
			MigrationName:       migrationName,
			ConcreteMigrationId: strconv.Itoa(int(sec)),
			Dependencies:        "",
			BlueprintName:       opts.Blueprint,
		}
		if lastMigrationId > 0 {
			concreteData.Dependencies = "{" + fmt.Sprintf(`"%s.%s"`, opts.Blueprint, strconv.Itoa(lastMigrationId)) + "}"
		}
		if err = concreteTpl.Execute(&concreteTplBuffer, concreteData); err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(pathToConcreteMigrationsFile, concreteTplBuffer.Bytes(), 0755)
		if err != nil {
			panic(err)
		}
		integrateMigrationIntoRegistryTpl := template.Must(template.New("integratemigrationintoregistry").Parse(initializeMigrationRegistryTpl))
		integrateMigrationIntoRegistryData := struct {
			MigrationName string
		}{
			MigrationName: migrationName,
		}
		var integrateMigrationIntoRegistryTplBuffer bytes.Buffer
		if err = integrateMigrationIntoRegistryTpl.Execute(&integrateMigrationIntoRegistryTplBuffer, integrateMigrationIntoRegistryData); err != nil {
			panic(err)
		}
		read, err := ioutil.ReadFile(pathToBaseMigrationsFile)
		if err != nil {
			panic(err)
		}
		newContents := strings.Replace(
			string(read),
			"// placeholder to insert next migration",
			integrateMigrationIntoRegistryTplBuffer.String()+"\n    // placeholder to insert next migration", -1)
		err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(newContents), 0755)
		if err != nil {
			panic(err)
		}
		fmt.Printf(
			"Created migration for blueprint %s with name %s\n",
			opts.Blueprint,
			opts.Message,
		)
	}
	return nil
}

func (command CreateMigration) GetHelpText() string {
	return "Create migration for your blueprint"
}

func ensureDatabaseIsReadyForMigrationsAndReadAllApplied() []Migration {
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	dbForMigrations := uadminDatabase.Db
	err := dbForMigrations.AutoMigrate(Migration{})
	if err != nil {
		panic(fmt.Errorf("error while preparing database for migrations: %s", err))
	}
	var appliedMigrations []Migration
	dbForMigrations.Find(&appliedMigrations)
	return appliedMigrations
}

type UpMigrationOptions struct {
}

type UpMigration struct {
}

func (command UpMigration) Proceed(subaction string, args []string) error {
	ensureDatabaseIsReadyForMigrationsAndReadAllApplied()
	for traverseMigrationResult := range appInstance.BlueprintRegistry.TraverseMigrations() {
		if traverseMigrationResult.Error != nil {
			panic(traverseMigrationResult.Error)
		}
		if traverseMigrationResult.Node.IsApplied() {
			continue
		}
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		appliedMigration := Migration{}
		uadminDatabase.Db.Where(
			&Migration{MigrationName: traverseMigrationResult.Node.GetMigration().GetName()},
		).First(&appliedMigration)
		if appliedMigration.ID != 0 {
			continue
		}
		color.Blue("Applying migration %s", traverseMigrationResult.Node.GetMigration().GetName())
		//if appInstance.Config.D.Db.Default.Type == "sqlite" {
		//	traverseMigrationResult.Node.Apply()
		//} else {
		db := uadminDatabase.Db
		db.Transaction(func(tx *gorm.DB) error {
			uadminDatabase1 := &core.UadminDatabase{
				Adapter: uadminDatabase.Adapter,
				Db:      tx,
			}
			res := traverseMigrationResult.Node.Apply(uadminDatabase1)
			if res == nil {
				tx.Create(
					&Migration{
						MigrationName: traverseMigrationResult.Node.GetMigration().GetName(),
						AppliedAt:     time.Now(),
					},
				)
			}
			return res
		})
		//}

	}
	return nil
}

func (command UpMigration) GetHelpText() string {
	return "Upgrade your database"
}

type DownMigrationOptions struct {
	MigrationName string `short:"m" required:"false" default:"" description:"Migration downgrade your database to"`
}

type DownMigration struct {
}

func (command DownMigration) Proceed(subaction string, args []string) error {
	var opts = &DownMigrationOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if err != nil {
		panic(err)
	}
	ensureDatabaseIsReadyForMigrationsAndReadAllApplied()
	for traverseMigrationResult := range appInstance.BlueprintRegistry.TraverseMigrationsDownTo(opts.MigrationName) {
		if traverseMigrationResult.Error != nil {
			panic(traverseMigrationResult.Error)
		}
		migrationName := traverseMigrationResult.Node.GetMigration().GetName()
		appliedMigration := Migration{}
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		result := uadminDatabase.Db.Where(
			"migration_name = ?", migrationName,
		).First(&appliedMigration)
		if result.RowsAffected == 0 {
			continue
			//panic(
			//	fmt.Sprintf(
			//		"Migration with name %s was not applied, so we can't downgrade database", migrationName,
			//	),
			//)
		}
		color.Blue("Downgrading migration %s", traverseMigrationResult.Node.GetMigration().GetName())
		//if appInstance.Config.D.Db.Default.Type == "sqlite" {
		//	uadminDatabase := interfaces.NewUadminDatabase()
		//	traverseMigrationResult.Node.Downgrade()
		//	defer uadminDatabase.Close()
		//	db := uadminDatabase.Db
		//	db.Unscoped().Delete(&appliedMigration)
		//} else {
		db := uadminDatabase.Db
		db.Transaction(func(tx *gorm.DB) error {
			uadminDatabase1 := &core.UadminDatabase{
				Adapter: uadminDatabase.Adapter,
				Db:      tx,
			}
			res := traverseMigrationResult.Node.Downgrade(uadminDatabase1)
			if res == nil {
				tx.Unscoped().Delete(&appliedMigration)
			}
			return res
		})
		//		}
	}
	return nil
}

func (command DownMigration) GetHelpText() string {
	return "Downgrade your database"
}

type DetermineConflictsMigration struct {
}

func (command DetermineConflictsMigration) Proceed(subaction string, args []string) error {
	ensureDatabaseIsReadyForMigrationsAndReadAllApplied()
	isEverythingOk := true
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	for traverseMigrationResult := range appInstance.BlueprintRegistry.TraverseMigrations() {
		if traverseMigrationResult.Error != nil {
			isEverythingOk = false
			core.Trail(core.WARNING, "Potential problems with migrations %s", traverseMigrationResult.Error.Error())
		}
		appliedMigration := Migration{}
		uadminDatabase.Db.Where(
			&Migration{MigrationName: traverseMigrationResult.Node.GetMigration().GetName()},
		).First(&appliedMigration)
		if appliedMigration.ID != 0 {
			continue
		}
		isEverythingOk = false
		core.Trail(core.WARNING, "Not applied migration: %s", traverseMigrationResult.Node.GetMigration().GetName())
	}
	if !isEverythingOk {
		return fmt.Errorf("determined some problems with migrations")
	}
	return nil
}

func (command DetermineConflictsMigration) GetHelpText() string {
	return "Determine if there any conflicts in migrations"
}
