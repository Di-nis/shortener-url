statictests:
	go vet -vettool=statictest ./...

autotests1:
	shortenertest -test.v -test.run=^TestIteration1$ \
      -binary-path=cmd/shortener/shortener

autotests11:
	shortenertest -test.v -test.run=^TestIteration11$ \
      -binary-path=cmd/shortener/shortener \
      -database-dsn=postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable

autotests9:
	shortenertest -test.v -test.run=^TestIteration9$ \
      -binary-path=cmd/shortener/shortener \
      -source-path=. \
      -file-storage-path=database.log

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

all: statictests autotests14