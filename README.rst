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
      You can set extra headers with the fourth argument as described in the HPCloud
      documentation
    */
    if err := acc.ObjectStoreUpload("/path/to/file", "container", "as", nil); err != nil {
        Log.Fatal(err)
    }
    
    if err := acc.ObjectStoreDelete("/path/to/file/on/objectstore/"); err != nil {
        Log.Fatal(err)
    }
    
    expires_utc := "2147483647"
    list_objects, err := range acc.ListObjects("/container/")
    if err != nil {
        Log.Fatal(err)
    }
    for _, entry :=  *list_objects {
         fmt.Println(acc.TemporaryURL(entry.Name, expires_utc))
    }
    
Any questions or bugs, please let me know and I will be happy to look over pull
requests or feature ideas.
