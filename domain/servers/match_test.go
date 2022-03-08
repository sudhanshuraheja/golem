package servers

import (
	"testing"
)

func TestMatcher(t *testing.T) {
	// servers := []Server{
	// 	{
	// 		Name:     "one",
	// 		HostName: []string{"one"},
	// 		Port:     22,
	// 		Tags:     []string{"one"},
	// 	},
	// 	{
	// 		Name:     "two",
	// 		PublicIP: localutils.StrPtr("127.0.0.1"),
	// 		HostName: []string{"two"},
	// 		Port:     22,
	// 		Tags:     []string{"one", "two"},
	// 	},
	// 	{
	// 		Name:     "three",
	// 		HostName: []string{"three"},
	// 		Port:     22,
	// 		Tags:     []string{"one", "two", "three"},
	// 	},
	// }

	// found, err := NewMatch("tags", "contains", "one").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 3, len(found))

	// found, err = NewMatch("tags", "contains", "two").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 2, len(found))

	// found, err = NewMatch("tags", "contains", "three").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 1, len(found))

	// found, err = NewMatch("tags", "not-contains", "two").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 1, len(found))

	// found, err = NewMatch("public_ip", "=", "127.0.0.1").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 1, len(found))

	// found, err = NewMatch("public_ip", "!=", "127.0.0.1").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 2, len(found))

	// found, err = NewMatch("port", "=", "22").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 3, len(found))

	// found, err = NewMatch("port", "!=", "22").Find(servers)
	// utils.Test().Nil(t, err)
	// utils.Test().Equals(t, 0, len(found))

}
