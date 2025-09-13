.PHONY: autotests

# test: run all autotests
autotests:
	shortenertest -test.v -test.run=^TestIteration12$ \
		-binary-path=cmd/shortener/shortener \
		-database-dsn=postgresql://postgres:postgres@localhost:5432/urls_db?sslmode=disable