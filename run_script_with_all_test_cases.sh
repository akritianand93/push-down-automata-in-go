echo "----------------Grammar 1----------------"

go run processor.go grammar1.json input_g1_0011.txt

go run processor.go grammar1.json input_g1_000011.txt

go run processor.go grammar1.json input_g1_1010.txt

go run processor.go grammar1.json input_g1_001111.txt

go run processor.go grammar1.json input_g1_00001111.txt

go run processor.go grammar1.json input_g1_ab.txt

go run processor.go grammar1.json input_g1_empty_string.txt

echo "----------------Grammar 2----------------"

go run processor.go grammar2.json input_g2_0111.txt

go run processor.go grammar2.json input_g2_011000.txt

go run processor.go grammar2.json input_g2_0101010.txt

echo "\n----------This script runs all sample test cases for grammar1 and grammar2-----------"
