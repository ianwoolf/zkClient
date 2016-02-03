## zkClient

zk client, zk demo


### usage

    ./zkClient -h
    Usage of ./zkClient:zkClient -c [param] -p [param] (-d -v -flags)
    command list: get/set/child/creat/watchExist/watchChildren/delChildNode

      -c string
    	    command, such as: get/set/child/creat/watchExist/watchChildren/del
      -d string
        	string data (default "test set2")
      -flags int
        	flag: 0-Permanent 1 2-sequence
      -p string
        	node path (default "/mynode/test")
      -v int
        	data version
    	
### command
    go build
    ./zkClient -c get -p /mynode/test
    ./zkClient -c child -p /mynode/test
    ./zkClient -c create -p /mynode/test
    ./zkClient -c child -p /mynode/test

    # first get，get node version, then set node value
    ./zkClient -c get -p /mynode/test/test1
    ./zkClient -c set -p /mynode/test/test1 -d "test test515" -v 20

    ./zkClient -c del -p /mynode/test # have no "/" at last
    