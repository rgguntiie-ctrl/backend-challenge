.PHONY: run.% build.% migrate.% tidy clean test test.%

# Default tidy
tidy:
	go mod tidy

# Dynamic run: make run.<app_name>
run.%:
	@echo "Running $* ..."
	go run ./cmd/$*

# Dynamic build: make build.<app_name>
build.%:
	@echo "Building $* ..."
	go build -o $* ./cmd/$*

# Dynamic migrate: make migrate.<app_name>
migrate.%:
	@echo "Migrating $* ..."
	go run ./cmd/$* migrate

# --- Unit test ---

# รันทุกไฟล์ unit test
test:
	@echo "Running all unit tests ..."
	go test -v ./unit_tests/...

# รันเฉพาะ module/ไฟล์
# ตัวอย่าง: make test.user
test.%:
	@echo "Running unit tests for $* ..."
	go test -v ./unit_tests/$*

