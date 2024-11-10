#!/bin/bash
green() {
  echo -e "\033[32m$1\033[0m"
}

red() {
  echo -e "\033[31m$1\033[0m"
}

cyan() {
  echo -e "\033[36m$1\033[0m"
}

cyan "🧪 Running Unit test..."
echo "----------------------------------------------------------------------------------"
if go test -race ./internal/service/... ./internal/handler/...; then
  echo "----------------------------------------------------------------------------------"
  green "🎉 All tests passed successfully! Great job! Keep up the excellent work! 💪\n"
else
  echo "----------------------------------------------------------------------------------"
  red "❌ Oops! Some tests failed. Let's fix those issues and try again. 📈\n"
  exit 1
fi
