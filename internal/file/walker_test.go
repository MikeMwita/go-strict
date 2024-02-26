package file_test

/* ------------------------------------ Imports -------------------------------- */

/* ------------------------------------ Methods/Functions -------------------------------- */

//func Test_Walker_Wrong_Directory(t *testing.T) {
//	rootDir := libFile.RootDir()
//
//	// create some files first
//	testDir := rootDir + "/file/temp"
//
//	files, err := libFile.Walker(testDir, "*.txt", 0)
//	assert.NoError(t, err)
//	assert.Len(t, files, 0)
//}
//
//func Test_Walker_Single_Directory_Wrong_Filter(t *testing.T) {
//	rootDir := libFile.RootDir()
//
//	// create some files first
//	testDir := rootDir + "/file/temp"
//	defer os.RemoveAll(testDir)
//
//	for i := 0; i < 10; i++ {
//		err := createTestFile(testDir)
//		assert.NoError(t, err)
//	}
//
//	files, err := libFile.Walker(testDir, "*.wrong", 0)
//	assert.NoError(t, err)
//	assert.Len(t, files, 0)
//}
//
//func Test_Walker_Single_Directory(t *testing.T) {
//	rootDir := libFile.RootDir()
//
//	// create some files first
//	testDir := rootDir + "/file/temp"
//	defer os.RemoveAll(testDir)
//
//	for i := 0; i < 10; i++ {
//		err := createTestFile(testDir)
//		assert.NoError(t, err)
//	}
//
//	// create some files in a subdirectory
//	testSubDir := rootDir + "/file/temp/subdir"
//	for i := 0; i < 5; i++ {
//		err := createTestFile(testSubDir)
//		assert.NoError(t, err)
//	}
//
//	files, err := libFile.Walker(testDir, "*.txt", 0)
//	assert.NoError(t, err)
//	assert.Len(t, files, 10)
//}
//
//func Test_Walker_SubDirectory(t *testing.T) {
//	rootDir := libFile.RootDir()
//
//	// create some files first
//	testDir := rootDir + "/file/temp"
//	defer os.RemoveAll(testDir)
//
//	for i := 0; i < 10; i++ {
//		err := createTestFile(testDir)
//		assert.NoError(t, err)
//	}
//
//	// create some files in a subdirectory
//	testSubDir := rootDir + "/file/temp/subdir"
//	for i := 0; i < 5; i++ {
//		err := createTestFile(testSubDir)
//		assert.NoError(t, err)
//	}
//
//	// create some files in a subdirectory
//	testDeeperSubDir := rootDir + "/file/temp/subdir/deeper/"
//	for i := 0; i < 2; i++ {
//		err := createTestFile(testDeeperSubDir)
//		assert.NoError(t, err)
//	}
//
//	files, err := libFile.Walker(testDir, "*.txt", 1)
//	assert.NoError(t, err)
//	assert.Len(t, files, 15)
//}
//
//func Test_Walker_AllDirectories(t *testing.T) {
//	rootDir := libFile.RootDir()
//
//	// create some files first
//	testDir := rootDir + "/file/temp"
//	defer os.RemoveAll(testDir)
//
//	for i := 0; i < 10; i++ {
//		err := createTestFile(testDir)
//		assert.NoError(t, err)
//	}
//
//	// create some files in a subdirectory
//	testSubDir := rootDir + "/file/temp/subdir"
//	for i := 0; i < 5; i++ {
//		err := createTestFile(testSubDir)
//		assert.NoError(t, err)
//	}
//
//	// create some files in a subdirectory
//	testDeeperSubDir := rootDir + "/file/temp/subdir/deeper/"
//	for i := 0; i < 2; i++ {
//		err := createTestFile(testDeeperSubDir)
//		assert.NoError(t, err)
//	}
//
//	files, err := libFile.Walker(testDir, "*.txt", -1)
//	assert.NoError(t, err)
//	assert.Len(t, files, 17)
//}
//
//// createTestFile creates a random test file.
//func createTestFile(testDir string) error {
//	rsuffix := randomString(5)
//	testFile := testDir + "/testfile_" + rsuffix + ".txt"
//	err := libFile.TextWriter("Marcello Holland", testFile)
//
//	return err
//}
//
//// Returns an int >= min, < max.
//func randomInt(min, max int) int {
//	return min + rand.Intn(max-min)
//}
//
//// Generate a random string of A-Z chars with len = l.
//func randomString(length int) string {
//	bytes := make([]byte, length)
//
//	for i := 0; i < length; i++ {
//		bytes[i] = byte(randomInt(65, 90))
//	}
//
//	return string(bytes)
//}
//
//
//func BenchmarkWalker(b *testing.B) {
//	rootDir, _ := RootDir()
//	// create some files first
//	testDir := rootDir + "/file/temp"
//
//	defer os.RemoveAll(testDir)
//
//	for i := 0; i < 10; i++ {
//		_ = createTestFile(testDir)
//	}
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ { // use b.N for looping
//		_, _ = Walker(testDir, "*.txt", 0)
//	}
//}
//
//// Helper functions
//func createTestFile(testDir string) error {
//	rsuffix := randomString(5)
//	testFile := testDir + "/testfile_" + rsuffix + ".txt"
//	err := TextWriter("MikeMwita", testFile)
//	return err
//}
//
//func randomString(length int) string {
//	bytes := make([]byte, length)
//
//	for i := 0; i < length; i++ {
//		bytes[i] = byte(randomInt(65, 90))
//	}
//
//	return string(bytes)
//}
//
//func randomInt(min, max int) int {
//	return min + rand.Intn(max-min)
//}

/*
Benchmark_Walker-8   	  187741	     21920 ns/op	    5760 B/op	      97 allocs/op
Benchmark_Walker-8   	   57690	     20103 ns/op	    5760 B/op	      97 allocs/op

Benchmark_Walker-8   	   59488	     20287 ns/op	    5760 B/op	      97 allocs/op

*/
