package main

import (
	"os"
)

const (
	wodbySSHKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDT5YlyihKjNncd7qApnB/VUwMOP18OqKPF9icq0XpCog2on9QsX2bNjtVjQ7/dAI7yQd4HtmfdrCZh3WVtCEeCiLmiMlQQ3Ac6Kc/K3VluANdP9ssn+N0+RsDp9VsrfEvBYvNjRpuiZcs9XcVCWTGYTug45PLfl63NqmZ9LU4wWAn+kt4TqRgSXuYV+LFWtRE/zpmxc6JUDZGWJGltH9Ly8Bf+OmkIFuZChTUzktMYFizBkmQz+uf+ru/BtBAmiqUTW2gaeORM5nTmsBwU2gndbPBBwFV5I1BM8kpmsFRH/ZLSQf6LySxPQguSglT2b/N9y25naUZbMtoGMChFscR63SdzGgE76DEuyGSW8qPwurYpR4H+ExqlettVy+OYskuQSg09ez4h4lcnV5Dmtxk8durK+CCqM6c4kIaKUbnOEud0ztJvQIKQRpMFIm09ySUB57Nf6kCNsT8Lje9vVcArfnmYOM4fhkOvqPQc2i27r76jxk0+rCUeUzXZDvCxS1A6fiWPyj9FA9BDNPl/vkCXCzO3Bz/jiqCDL/Z5FzYiXEmcaOfVQ/pnGL2wTZdHCjXqS4AaxKNF/YrQ9HR3eUaml35aHxnvcrXIQpSkJRU3owWFr12JWllMi21SSRLamT7zcp3e4nWowj6jQBE9Wh+yhwsewxnXzk/MSjNbHrgOCQ==`
)

func addSSHKey() error {
	f, err := os.OpenFile("/root/.ssh/authorized_keys", os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(wodbySSHKey)
	if err != nil {
		return err
	}
	return nil
}