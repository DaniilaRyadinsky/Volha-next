package postgres

import (
	"net"
	"net/url"
	"strconv"
)

func URL(host string, port int, database string, user string, password string, sslMode string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/" + database,
	}
	if escaped := "/" + url.PathEscape(database); escaped != u.Path {
		u.RawPath = escaped
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}
