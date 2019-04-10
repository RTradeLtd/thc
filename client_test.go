package thc

import (
	"fmt"
	"os"
	"testing"
)

func Test_Call_String(t *testing.T) {
	checks := []struct {
		call Call
		want string
	}{
		{Login, "/auth/login"},
		{FileAddPublic, "/ipfs/public/file/add"},
		{PinAddPublic, "/ipfs/public/pin/%s"},
	}
	for _, check := range checks {
		if check.call.String() != check.want {
			t.Fatalf("received %s, expected %s", check.call.String(), check.want)
		}
	}
}

func Test_Call_FillParams_Pin_Add_Public(t *testing.T) {
	if PinAddPublic.FillParams("hello") != "/ipfs/public/pin/hello" {
		t.Fatal("bad fill params")
	}
	c := Call("/a/%s/c/%s")
	if c.FillParams("b", "d") != "/a/b/c/d" {
		t.Fatal("bad fill params")
	}
}

func Test_New_V2(t *testing.T) {
	type args struct {
		user, pass, url string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Dev", args{"testuser", "admin", DevURL}},
		{"Prod", args{"testuser", "admin", ProdURL}},
	}
	for _, tt := range tests {
		v2 := NewV2(tt.args.user, tt.args.pass, tt.args.url)
		if v2.auth.user != tt.args.user {
			t.Fatal("bad user set")
		}
		if v2.auth.pass != tt.args.pass {
			t.Fatal("bad path set")
		}
		if v2.url != tt.args.url {
			t.Fatal("bad url set")
		}
	}
}

func Test_Login(t *testing.T) {
	t.Skip("skipping integration test")
	v2 := NewV2(os.Getenv("USER"), os.Getenv("PASS"), DevURL)
	if err := v2.Login(); err != nil {
		t.Fatal(err)
	}
}

func Test_File_Add(t *testing.T) {
	t.Skip("skipping integration test")
	v2 := NewV2(os.Getenv("USER"), os.Getenv("PASS"), DevURL)
	if err := v2.Login(); err != nil {
		t.Fatal(err)
	}
	if _, err := v2.FileAdd("../../README.md", FileAddOpts{HoldTime: "5"}); err != nil {
		t.Fatal(err)
	}
}

func Test_Pin_Add(t *testing.T) {
	t.Skip("skipping integration test")
	v2 := NewV2(os.Getenv("USER"), os.Getenv("PASS"), DevURL)
	if err := v2.Login(); err != nil {
		t.Fatal(err)
	}
	if _, err := v2.PinAdd("QmY8VGk1QRd7ko87wk3YscWBRvokzDeH4xobJudCbGNM6B", "5"); err != nil {
		t.Fatal(err)
	}
}

func Test_Lens_Search(t *testing.T) {
	t.Skip()
	v2 := NewV2(os.Getenv("USER"), os.Getenv("PASS"), DevURL)
	if err := v2.Login(); err != nil {
		t.Fatal(err)
	}
	if resp, err := v2.SearchLens("blockchain"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(resp)
	}
}
