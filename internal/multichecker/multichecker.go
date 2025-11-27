// Package multichecker представляет мультичекер статических анализаторов кода.
// Целью набора анализаторов является предоставление максимально полезных,
// практичных диагностик при явном контроле над тем, какие проверки
// включены.
//
// Мультичекер включает следующие анализаторы:
// 1) Анализаторы `passes`:
// - asmdecl — проверяет корректность объявлений функций в asm.
// - assign — находит неправильные или подозрительные операции присваивания.
// - atomic — предупреждает о неправильном использовании sync/atomic.
// - bools — ищет примитивные логические выражения вида `x == true`.
// - buildtag — проверяет формат и корректность build-тегов.
// - cgocall — обнаруживает потенциально опасные вызовы C-кода.
// - composite — подсвечивает ошибки в литералах составных типов.
// - copylock — проверяет копирование структур, содержащих sync.Mutex/Once.
// - ctrlflow — анализирует вероятные ошибки в управлении потоком.
// - deepequalerrors — предупреждает о сравнении ошибок через reflect.DeepEqual.
// - errorsas — проверяет правильность использования errors.As.
// - httpresponse — ищет ошибки при работе с http.Response (например, забытый Body.Close).
// - ifaceassert — проверяет приведение типов к интерфейсам.
// - loopclosure — подсвечивает замыкания на переменные цикла.
// - lostcancel — предупреждает о потерянных вызовах Cancel() у контекста.
// - nilfunc — находит вызовы nil-функций.
// - nilness — анализирует возможные nil-deref ошибки.
// - printf — проверяет корректность форматных строк.
// - shift — предупреждает о неверных сдвигах битов.
// - sigchanyzer — анализирует операции с каналами сигналов os/signal.
// - sortslice — проверяет корректность работы с sort.Slice.
// - stdmethods — ищет неправильные реализации стандартных методов (Error, String и др.).
// - stringintconv — подсвечивает ошибки строковых и целочисленных преобразований.
// - structtag — проверяет теги структур.
// - testinggoroutine — обнаруживает вызовы t.Fatal внутри горутин.
// - tests — проверяет корректность тестовых функций.
// - unmarshal — подсвечивает возможные ошибки при unmarshal JSON/XML.
// - unreachable — ищет недостижимый код.
// - unsafeptr — предупреждает о неправильном использовании unsafe.Pointer.
// - unusedresult — ищет игнорирование значимых возвращаемых значений.
//
// 2) Анализаторы `stylecheck`:
// - stylecheck/st1000 — проверяет наличие и корректность комментария к пакету
// - stylecheck/st1020 — требует, чтобы комментарии к экспортируемым идентификаторам начинались с большой буквы
//
// 3) Анализатор класса `simple`:
// - simple/s1000 — упрощает избыточные конструкции, предлагая более простой и идиоматичный вариант
//
// 4) Анализаторы класса SA (staticcheck):
// - SA1000–SA1032 — ошибки использования стандартной библиотеки (неверные аргументы функций, неправильное применение API, устаревшие вызовы).
// - SA2000–SA2003 — корректность работы с goroutine и sync (неправильная передача WaitGroup, утечки горутин, проблемы с мьютексами).
// - SA3000–SA3001 — ошибки тестирования (неправильное использование t.FailNow/t.Parallel и другие проблемы тестов).
// - SA4000–SA4029 — логические ошибки и подозрительные конструкции (бессмысленные сравнения, всегда истинные/ложные условия, дублирующий код).
// - SA4030–SA4032 — неэффективные или некорректные операции над слайсами и map.
// - SA5000–SA5012 — ошибки в обработке ошибок (пропуск проверки ошибок, неправильно сформированные ошибки и др.).
// - SA6000–SA6006 — ошибки в работе с каналами (закрытие открытого канала, чтение/запись в неверных местах и др.).
// - SA9001–SA9009 — разнородные ошибки качества кода (неиспользуемые значения, подозрительные конструкции, потенциальные баги).
//
// 5) Дополнительный пользовательский анализатор:
// - ErrExitCheckAnalyzer — запрещает прямой вызов os.Exit внутри функции
// main пакета main.
//
// Использование
//
// Минимальный пример главного файла multichecker:
//
// package main
//
// import (
// "golang.org/x/tools/go/analysis/multichecker"
// "example.com/your/module/multichecker"
// )
//
// func main() {
// multichecker.Main(multichecker.Analyzers)
// }
//
// Соберите бинарь и запускайте:
// ./multichecker ./...
// Анализаторы поддерживают стандартные флаги фреймворка analysis
// (например, `-v` для подробного вывода).
//
// Расширение и настройка
//
// - Для подавления диагностики используйте стандартные комментарии вида
// `//lint:ignore` или аналоги, если они поддерживаются анализатором.
//
//
// Отладка и советы
//
// - При внутренних ошибках анализаторов запустите multichecker с флагом `-v`
// для получения расширенной информации.
// - При дублирующихся предупреждениях проверьте, не сообщают ли об одном и том
// же разные анализаторы.
// - Рекомендуется фиксировать версии модулей staticcheck и x/tools в go.mod.
//
// Дальнейшая разработка
//
// - В проект можно добавлять собственные анализаторы, используя фреймворк
// go/analysis.
// - Тесты multichecker рекомендуется строить на основе примеров кода,
// проверяющих регистрацию всех анализаторов и корректность работы
// ErrExitCheckAnalyzer.
package multichecker

import (
	"strings"

	"github.com/Di-nis/shortener-url/internal/errexitchecker"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/simple/s1000"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1000"
	"honnef.co/go/tools/stylecheck/st1020"
	"honnef.co/go/tools/unused"

	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
)

// Run запускает мультичекер.
func Run() {
	analyzersPasses := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilness.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	var (
		checksNameSA = "SA"
		checksAll    []*analysis.Analyzer
	)

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, checksNameSA) {
			checksAll = append(checksAll, v.Analyzer)
		}
	}
	checksAll = append(checksAll, st1000.SCAnalyzer.Analyzer, st1020.SCAnalyzer.Analyzer)
	checksAll = append(checksAll, analyzersPasses...)
	checksAll = append(checksAll, errexitchecker.ErrExitCheckAnalyzer)
	checksAll = append(checksAll, unused.Analyzer.Analyzer, s1000.Analyzer)

	multichecker.Main(
		checksAll...,
	)
}
