package hpcloud

func (a Access) ObjectStoreUpload(filename, container, as string, header *http.Header) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s/%s", OBJECT_STORE, a.TenantID, container, as)
	fmt.Println(path)
	req, err := http.NewRequest("PUT", path, f)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.Token())
	if header != nil {
		for key, value := range header {
			for _, s := range value {
				req.Header.Add(key, s)
			}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
