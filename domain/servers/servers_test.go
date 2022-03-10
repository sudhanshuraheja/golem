package servers

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

func TestDomainIPs(t *testing.T) {
	d := NewDomainIP("localhost", "127.0.0.1")
	ds := DomainIPs{}
	ds.Append(*d)
	utils.Test().Equals(t, 1, len(ds))

	ds2 := DomainIPs{}
	ds2.Append(*d)
	utils.Test().Equals(t, 1, len(ds2))

	ds.Merge(ds2)
	utils.Test().Equals(t, 2, len(ds))
}

func TestServer(t *testing.T) {
	log := logger.NewCLILogger(6, 8)
	s1 := Server{
		Name:      "one",
		HostName:  &[]string{"one"},
		PublicIP:  localutils.StrPtr("1.2.3.4"),
		PrivateIP: localutils.StrPtr("1.2.3.5"),
		Port:      22,
		Tags:      &[]string{"one"},
	}

	s1.Display(log, "")
	s1.Display(log, "one")
	s1.Display(log, "two")

	host, err := s1.GetHostName()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "1.2.3.4", host)

	s1.PublicIP = nil
	host, err = s1.GetHostName()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "one", host)

	s1.HostName = nil
	host, err = s1.GetHostName()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "1.2.3.5", host)

	s1.PrivateIP = nil
	_, err = s1.GetHostName()
	utils.Test().Contains(t, err.Error(), "could not find")

	s1.PublicIP = localutils.StrPtr("1.2.3.4")
	s1.HostName = &[]string{"one"}
	s1.PrivateIP = localutils.StrPtr("1.2.3.5")

	found, err := s1.Search(*NewMatch("name", "=", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	found, err = s1.Search(*NewMatch("public_ip", "=", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)

	found, err = s1.Search(*NewMatch("private_ip", "=", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)

	found, err = s1.Search(*NewMatch("hostname", "contains", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	found, err = s1.Search(*NewMatch("user", "=", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)

	found, err = s1.Search(*NewMatch("port", "=", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)

	found, err = s1.Search(*NewMatch("tags", "contains", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	_, err = s1.Search(*NewMatch("tag", "=", "one"))
	utils.Test().Contains(t, err.Error(), "does not support attribute")

	s1.Name = ""
	found, err = s1.Search(*NewMatch("tags", "contains", "one"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)
}

func TestServers(t *testing.T) {
	log := logger.NewCLILogger(6, 8)
	s1 := Server{
		Name:      "one",
		HostName:  &[]string{"one"},
		PublicIP:  localutils.StrPtr("1.2.3.4"),
		PrivateIP: localutils.StrPtr("1.2.3.5"),
		Port:      22,
		Tags:      &[]string{"one"},
	}
	s2 := Server{
		Name:     "two",
		PublicIP: localutils.StrPtr("127.0.0.1"),
		HostName: &[]string{"two"},
		Port:     22,
		Tags:     &[]string{"one", "two"},
	}
	s3 := Server{
		Name:     "three",
		HostName: &[]string{"three"},
		Port:     22,
		Tags:     &[]string{"one", "two", "three"},
	}

	ss1 := Servers{}
	ss1.Append(s1)
	utils.Test().Equals(t, 1, len(ss1))

	ss2 := Servers{}
	ss2.Append(s2)
	ss2.Append(s3)
	utils.Test().Equals(t, 2, len(ss2))

	ss1.Merge(ss2)
	utils.Test().Equals(t, 3, len(ss1))

	ss1.Display(log, "")

	d := NewDomainIP("localhost", "127.0.0.1")
	ds := DomainIPs{}
	ds.Append(*d)

	ss1.MergesHosts(&ds)
	utils.Test().Equals(t, 2, len(*ss1[1].HostName))

	ss3, err := ss1.Search(*NewMatch("name", "=", "two"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "two", ss3[0].Name)

	_, err = ss1.Search(*NewMatch("tag", "=", "two"))
	utils.Test().Contains(t, err.Error(), "does not support attribute")
}
