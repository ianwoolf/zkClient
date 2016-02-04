## zkClient

zk client, zk demo


### usage

    ./zkClient -h
    Usage of ./zkClient:zkClient -c [param] -p [param] (-d -v -flags)
    command list: get/set/child/creat/watchExist/watchChildren/delChildNode

      -s value
          zk server list. e.g: 127.0.0.1:2181,127.0.0.2:2181. 
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
    ./zkClient -s 127.0.0.1:2181 -c get -p /mynode/test
    ./zkClient -s 127.0.0.1:2181 -c child -p /mynode/test
    # flags: 0-Permanent 2-sequence
    ./zkClient -s 127.0.0.1:2181 -c create -p /mynode/test/test6 -d "content" -flags 0
    ./zkClient -s 127.0.0.1:2181 -c create -p /mynode/test/se-job -d "content" -flags 2

    # first get，get node version, then set node value
    ./zkClient -s 127.0.0.1:2181 -c get -p /mynode/test/test1
    ./zkClient -s 127.0.0.1:2181 -c set -p /mynode/test/test1 -d "test test515" -v 20

    ./zkClient -s 127.0.0.1:2181 -c del -p /mynode/test # have no "/" at last

    # watch todo: watch by param type: exist/child
    ./zkClient -s 127.0.0.1:2181 -c watch -p /mynode/test
    ./zkClient -s 127.0.0.1:2181 -c watch -p /mynode/test/test6
    