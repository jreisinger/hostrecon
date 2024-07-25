package geo

import (
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jreisinger/hostrecon"
	"github.com/oschwald/geoip2-golang"
)

type DbIpCom struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
}

func (DbIpCom) Recon(target string) hostrecon.Info {
	info := hostrecon.Info{Host: target, Kind: "db-ip.com"}

	basename := fmt.Sprintf("dbip-city-lite-%s.mmdb", time.Now().Format("2006-01"))
	url := fmt.Sprintf("https://download.db-ip.com/free/%s.gz", basename)
	path := filepath.Join("/tmp", basename)

	f, err := os.Stat(path)
	if (err != nil && os.IsNotExist(err)) || isOlderThanOneWeek(f.ModTime()) {
		body, err := download(url)
		if err != nil {
			info.Err = err
			return info
		}
		if err := extract(body, path); err != nil {
			info.Err = err
			return info
		}
	} else if err != nil {
		info.Err = err
		return info
	}

	db, err := geoip2.Open(path)
	if err != nil {
		info.Err = err
		return info
	}
	defer db.Close()

	addrs, err := net.LookupHost(target)
	if err != nil {
		info.Err = err
		return info
	}

	for _, addr := range addrs {
		geo, err := db.City(net.ParseIP(addr))
		if err != nil {
			info.Err = err
			return info
		}
		location := fmt.Sprintf("%s: %s %s", addr, geo.City.Names["en"], geo.Country.IsoCode)
		info.Info = append(info.Info, location)
	}

	return info
}

func isOlderThanOneWeek(t time.Time) bool {
	return time.Since(t) > 7*24*time.Hour
}

func download(url string) (r io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading %v: %v", url, resp.Status)
	}
	return resp.Body, nil
}

func extract(r io.ReadCloser, filename string) error {
	defer r.Close()

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, gzipReader); err != nil {
		return err
	}

	return nil
}
