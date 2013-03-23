.. image:: https://travis-ci.org/AeroNotix/hpcloud.png?branch=master

HPCloud Go bindings
===================


These bindings are incredibly new and are definitely not complete, however I
am certainly looking for help!

A few examples:

.. code-block:: go

	acc, err := hpcloud.Authenticate(username, password, tenantID)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	  Upload files easily to the object store, their metadata will be set
	  appropriately. The file will be MD5'd for end-to-end
	  integrity checks.

	  You can set extra headers with the fourth argument as described in the HPCloud
	  documentation.
	*/
	if err := acc.ObjectStoreUpload("/path/to/file", "container", nil); err != nil {
		Log.Fatal(err)
	}

	/* Delete items */
	if err := acc.ObjectStoreDelete("/path/to/file/on/objectstore/"); err != nil {
		Log.Fatal(err)
	}

	/* List objects in containers */
	expires_utc := "2147483647"
	list_objects, err := range acc.ListObjects("/container/")
	if err != nil {
		Log.Fatal(err)
	}
	for _, entry :=  *list_objects {
		 fmt.Println(acc.TemporaryURL(entry.Name, expires_utc))
	}

	/* Create new servers, easily */
	s, err := acc.CreateServer(hpcloud.Server{
		FlavorRef: hpcloud.XSmall,
		Name:      "MyAwesomeNewServer",
		Key:       "me",
		ImageRef:  hpcloud.DebianSqueeze6_0_3Server,
	})
	if err != nil {
		Log.Fatal(err)
	}
	fmt.Sprintf("Status: %s\nID: %d\n", s.S.Status, s.S.ID)

	/* Delete that server we just created */
	fmt.Println(acc.DeleteServer(s.S.ID))


Any questions or bugs, please let me know and I will be happy to look over pull
requests or feature ideas.
