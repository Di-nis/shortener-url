statictests:
	go vet -vettool=statictest ./...

autotests1:
	shortenertest -test.v -test.run=^TestIteration1$ \
      -binary-path=cmd/shortener/shortener

autotests9:
	shortenertest -test.v -test.run=^TestIteration9$ \
      -binary-path=cmd/shortener/shortener \
      -source-path=. \
      -file-storage-path=database.log

autotests10:
	shortenertest -test.v -test.run=^TestIteration10$ \
      -binary-path=cmd/shortener/shortener \
      -source-path=. \
      -database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests11:
	shortenertest -test.v -test.run=^TestIteration11$ \
      -binary-path=cmd/shortener/shortener \
      -database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests12:
	shortenertest -test.v -test.run=^TestIteration12$ \
      -binary-path=cmd/shortener/shortener \
      -database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests13:
	shortenertest -test.v -test.run=^TestIteration13$ \
		-binary-path=cmd/shortener/shortener \
		-database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests14:
	shortenertestbeta -test.v -test.run=^TestIteration14$ \
		-binary-path=cmd/shortener/shortener \
		-database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests15:
	shortenertestbeta -test.v -test.run=^TestIteration15$ \
		-binary-path=cmd/shortener/shortener \
		-database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable \
		-source-path=.

autotests17:
	shortenertest -test.v -test.run=^TestIteration17$ \
		-source-path=. \

autotests18:
	shortenertest -test.v -test.run=^TestIteration18$ \
		-source-path=. \

all:  statictests autotests1 autotests9 autotests10 autotests11 autotests12 autotests13 autotests14 autotests15