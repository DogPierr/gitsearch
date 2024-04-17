package gitsearch

import (
	"fmt"

	"github.com/DogPierr/gitsearch/gitsearch/internal/dataparser"
	"github.com/spf13/cobra"
)

var (
	repository   string
	revision     string
	orderby      string
	usecommitter bool
	format       string
	extensions   string
	languages    string
	exclude      string
	restrictto   string
)

func run(cmd *cobra.Command, args []string) {
	err := dataparser.Init(
		repository,
		revision,
		orderby,
		usecommitter,
		format,
		extensions,
		languages,
		exclude,
		restrictto,
	)
	if err != nil {
		panic(err)
	}
	result, err := dataparser.ShowData()
	if err != nil {
		panic(err)
	}
	fmt.Print(result)
}

var gitsearch = &cobra.Command{
	Use:   "gitsearch",
	Short: "Подсчёта статистик авторов git репозитория",
	Long:  "Подсчёта статистик авторов git репозитория. Все статистики считаются для состояния репозитория на момент конкретного коммита.",
	Run:   run,
}

func Init() int {
	gitsearch.PersistentFlags().StringVarP(&repository, "repository", "r", ".", "путь до Git репозитория; по умолчанию текущая директория")
	gitsearch.PersistentFlags().StringVarP(&revision, "revision", "v", "HEAD", "указатель на коммит; HEAD по умолчанию")
	gitsearch.PersistentFlags().StringVarP(&orderby, "order-by", "o", "lines", "ключ сортировки результатов; один из lines (дефолт), commits, files")
	gitsearch.PersistentFlags().BoolVarP(&usecommitter, "use-committer", "c", false, "булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера")
	gitsearch.PersistentFlags().StringVarP(&format, "format", "f", "tabular", "формат вывода; один из tabular (дефолт), csv, json, json-lines")
	gitsearch.PersistentFlags().StringVarP(&extensions, "extensions", "e", "", "список расширений, сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например, '.go,.md'")
	gitsearch.PersistentFlags().StringVarP(&languages, "languages", "l", "", "список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например 'go,markdown'")
	gitsearch.PersistentFlags().StringVarP(&exclude, "exclude", "d", "", "набор Glob паттернов, исключающих файлы из расчёта, например 'foo/*,bar/*'")
	gitsearch.PersistentFlags().StringVarP(&restrictto, "restrict-to", "x", "", "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")

	if err := gitsearch.Execute(); err != nil {
		return 1
	}
	return 0
}
