test:
	go test -v -run='^\QTest_Check_' ./base
	go test -v -run='^\QTest_Check_' ./goTicker
race:
	go test -race -run='^\QTest_Check_' ./base
	go test -race -run='^\QTest_Check_' ./goTicker