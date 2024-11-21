this is an attempt to solve the 1 billion row challenge in go

you can read more about it here https://www.morling.dev/blog/one-billion-row-challenge/

here is a more advanced solution written in go using custom hashmaps and stuffs https://benhoyt.com/writings/go-1brc/
this codebase is made by reading (blatantly copying) the implementation details from the above solution

current status:
 proccessing takes:
   - about 9sec on a ryzen 5 5600 system with 16gb ram running ubuntu


to run this solution:
  - generate the measurements.txt file by following instructions from the 1brc repo (here https://github.com/gunnarmorling/1brc ) <i> if you managed to do that you can solve the solution already don't bother running my code </i>
  - edit the dataFilePath variable in main.go to your file location,
  - run with `go run main.go`
