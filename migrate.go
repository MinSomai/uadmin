package uadmin

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/uadmin/uadmin/interfaces"
	"gorm.io/gorm"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	gorm.Model
	MigrationName string `gorm:"index:migration_migration_name,unique"`
	AppliedAt  time.Time
}

type MigrateCommand struct {

}
func (c MigrateCommand) Proceed(subaction string, args []string) error {
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}
	createCommand := new(CreateMigration)

	commandRegistry.addAction("create", interfaces.ICommand(createCommand))
	upCommand := new(UpMigration)

	commandRegistry.addAction("up", interfaces.ICommand(upCommand))
	downCommand := new(DownMigration)

	commandRegistry.addAction("down", interfaces.ICommand(downCommand))
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

var re = regexp.MustCompile("[[:^ascii:]]")

func prepareMigrationName(message string) string {
	now := time.Now()
	sec := now.Unix()
	message = re.ReplaceAllLiteralString(message, "")
	if len(message) > 10 {
		message = message[:10]
	}
	return fmt.Sprintf("%s_%d", message, sec)
}

type CreateMigrationOptions struct {
	Message string `short:"m" required:"true" description:"Describe what is this migration for"`
	Blueprint string `short:"b" required:"true" description:"Blueprint you'd like to create migration for'"`
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
`
		fmt.Printf(help)
		return nil
	}
	if err != nil {
		panic(err)
	}
	const concreteMigrationTpl = `
package migrations

import (
    "github.com/uadmin/uadmin/utils"
)

type {{.MigrationName}} struct {
}

func (m {{.MigrationName}}) GetName() string {
    return "{{.BlueprintName}}.{{.ConcreteMigrationId}}"
}

func (m {{.MigrationName}}) GetId() int64 {
    return {{.ConcreteMigrationId}}
}

func (m {{.MigrationName}}) Up() {
}

func (m {{.MigrationName}}) Down() {
}

func (m {{.MigrationName}}) Deps() []string {
{{if .DependencyId}}    return []string{"{{.BlueprintName}}.{{.DependencyId}}"}{{else}}    return make([]string, 0){{end}}
}
`
	const initializeMigrationRegistryTpl = `
    BMigrationRegistry.AddMigration({{.MigrationName}}{})`
	const migrationRegistryCreationTpl = `
package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    // placeholder to insert next migration
}
`
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
	concreteData := struct{
		MigrationName string
		ConcreteMigrationId string
		DependencyId string
		BlueprintName string
	}{
		MigrationName: migrationName,
		ConcreteMigrationId: strconv.Itoa(int(sec)),
		DependencyId: "",
		BlueprintName: opts.Blueprint,
	}
	if lastMigrationId > 0 {
		concreteData.DependencyId = strconv.Itoa(lastMigrationId)
	}
	if err = concreteTpl.Execute(&concreteTplBuffer, concreteData); err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pathToConcreteMigrationsFile, concreteTplBuffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
	integrateMigrationIntoRegistryTpl := template.Must(template.New("integratemigrationintoregistry").Parse(initializeMigrationRegistryTpl))
	integrateMigrationIntoRegistryData := struct{
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
		integrateMigrationIntoRegistryTplBuffer.String() + "\n    // placeholder to insert next migration", -1)
	err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(newContents), 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Created migration for blueprint %s with name %s\n",
		opts.Blueprint,
		opts.Message,
	)
	return nil
}

func (command CreateMigration) GetHelpText() string {
	return "Create migration for your blueprint"
}

func ensureDatabaseIsReadyForMigrationsAndReadAllApplied() []Migration {
	err := appInstance.Database.ConnectTo("default").AutoMigrate(Migration{})
	if err != nil {
		panic(fmt.Errorf("error while preparing database for migrations: %s", err))
	}
	var appliedMigrations []Migration
	appInstance.Database.ConnectTo("default").Find(&appliedMigrations)
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
		appInstance.Database.ConnectTo("default").Create(
			&Migration{
				MigrationName: traverseMigrationResult.Node.GetMigration().GetName(),
				AppliedAt: time.Now(),
			},
		)
		traverseMigrationResult.Node.Apply()
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
		result := appInstance.Database.ConnectTo("default").Where(
			"migration_name = ?", migrationName,
		).First(&appliedMigration)
		if result.RowsAffected == 0 {
			panic(
				fmt.Sprintf(
					"Migration with name %s was not applied, so we can't downgrade database", migrationName,
				),
			)
		}
		traverseMigrationResult.Node.Downgrade()
		appInstance.Database.ConnectTo("default").Delete(&appliedMigration)
	}
	return nil
}

func (command DownMigration) GetHelpText() string {
	return "Downgrade your database"
}
