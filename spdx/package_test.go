package spdx

import "testing"

func TestPackage2_1CanGetVerificationCode(t *testing.T) {
	files := []*File2_1{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			Checksums:          []Checksum{{Value: "3333333333bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums:          []Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			Checksums:          []Checksum{{Value: "2222222222bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums:          []Checksum{{Value: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa", Algorithm: SHA1}},
		},
	}

	wantCode := PackageVerificationCode{Value: "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"}

	gotCode, err := MakePackageVerificationCode2_1(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_1CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := []*File2_1{
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			Checksums:          []Checksum{{Value: "3333333333bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			Checksums:          []Checksum{{Value: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa", Algorithm: SHA1}},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			Checksums:          []Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums:          []Checksum{{Value: "2222222222bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
	}

	wantCode := PackageVerificationCode{Value: "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"}

	gotCode, err := MakePackageVerificationCode2_1(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}
}

func TestPackage2_1GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := []*File2_1{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
		nil,
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums:          []Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: SHA1}},
		},
	}

	_, err := MakePackageVerificationCode2_1(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestPackage2_2CanGetVerificationCode(t *testing.T) {
	files := []*File2_2{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
	}

	wantCode := PackageVerificationCode{Value: "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"}

	gotCode, err := MakePackageVerificationCode2_2(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_2CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := []*File2_2{
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
	}

	wantCode := PackageVerificationCode{Value: "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"}

	gotCode, err := MakePackageVerificationCode2_2(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}
}

func TestPackage2_2GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := []*File2_2{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		nil,
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums: []Checksum{
				{
					Algorithm: SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
	}

	_, err := MakePackageVerificationCode2_2(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
