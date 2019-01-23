package main

//TODO : fix tests
//
//
//func TestGrpcApi(t *testing.T) {
//
//	port1, err := node.GetUnboundedPort()
//	port2, err := node.GetUnboundedPort()
//	assert.NoError(t, err, "Should be able to establish a connection on a port")
//
//	config.ConfigValues.JSONServerPort = port1
//	config.ConfigValues.GrpcServerPort = port2
//
//	const message = "Hello World"
//
//	grpcService := NewGrpcService(&net, ap)
//	grpcStatus := make(chan bool, 2)
//
//	// start a server
//	grpcService.StartService(grpcStatus)
//	<-grpcStatus
//
//	// start a client
//	addr := "localhost:" + strconv.Itoa(int(config.ConfigValues.GrpcServerPort))
//
//	// Set up a connection to the server.
//	conn, err := grpc.Dial(addr, grpc.WithInsecure())
//	if err != nil {
//		t.Fatalf("did not connect. %v", err)
//	}
//	defer conn.Close()
//	c := pb.NewSpaceMeshServiceClient(conn)
//
//	// call echo and validate result
//	r, err := c.Echo(context.Background(), &pb.SimpleMessage{Value: message})
//	if err != nil {
//		t.Fatalf("could not greet. %v", err)
//	}
//
//	assert.Equal(t, message, r.Value, "Expected message to be echoed")
//
//	// stop the server
//	grpcService.StopService()
//	<-grpcStatus
//}
//
//func TestJsonApi(t *testing.T) {
//
//	port1, err := node.GetUnboundedPort()
//	port2, err := node.GetUnboundedPort()
//	assert.NoError(t, err, "Should be able to establish a connection on a port")
//
//	config.ConfigValues.JSONServerPort = port1
//	config.ConfigValues.GrpcServerPort = port2
//
//	grpcService := NewGrpcService(&net, ap)
//	jsonService := NewJSONHTTPServer()
//
//	jsonStatus := make(chan bool, 2)
//	grpcStatus := make(chan bool, 2)
//
//	// start grp and json server
//	grpcService.StartService(grpcStatus)
//	<-grpcStatus
//
//	jsonService.StartService(jsonStatus)
//	<-jsonStatus
//
//	const message = "hello world!"
//	const contentType = "application/json"
//
//	// generate request payload (api input params)
//	reqParams := pb.SimpleMessage{Value: message}
//	var m jsonpb.Marshaler
//	payload, err := m.MarshalToString(&reqParams)
//	assert.NoError(t, err, "failed to marshal to string")
//
//	// Without this running this on Travis CI might generate a connection refused error
//	// because the server may not be ready to accept connections just yet.
//	time.Sleep(3 * time.Second)
//
//	url := fmt.Sprintf("http://127.0.0.1:%d/v1/example/echo", config.ConfigValues.JSONServerPort)
//	resp, err := http.Post(url, contentType, strings.NewReader(payload))
//	assert.NoError(t, err, "failed to http post to api endpoint")
//
//	defer resp.Body.Close()
//	buf, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	if got, want := resp.StatusCode, http.StatusOK; got != want {
//		t.Errorf("resp.StatusCode = %d; want %d", got, want)
//	}
//
//	var msg pb.SimpleMessage
//	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
//		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
//		return
//	}
//
//	if got, want := msg.Value, message; got != want {
//		t.Errorf("msg.Value = %q; want %q", got, want)
//	}
//
//	if value := resp.Header.Get("Content-Type"); value != contentType {
//		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
//	}
//
//	// stop the services
//	jsonService.StopService()
//	<-jsonStatus
//	grpcService.StopService()
//	<-grpcStatus
//}
//
//
//func TestJsonWalletApi_Errors(t *testing.T) {
//
//	port1, err := node.GetUnboundedPort()
//	port2, err := node.GetUnboundedPort()
//	assert.NoError(t, err, "Should be able to establish a connection on a port")
//	addrBytes := []byte{0x01}
//	config.ConfigValues.JSONServerPort = port1
//	config.ConfigValues.GrpcServerPort = port2
//
//	grpcService := NewGrpcService()
//	jsonService := NewJSONHTTPServer()
//
//	jsonStatus := make(chan bool, 2)
//	grpcStatus := make(chan bool, 2)
//
//	// start grp and json server
//	grpcService.StartService(grpcStatus)
//	<-grpcStatus
//
//	jsonService.StartService(jsonStatus)
//	<-jsonStatus
//
//	const contentType = "application/json"
//
//	// generate request payload (api input params)
//	reqParams := pb.AccountId{Address: addrBytes}
//	var m jsonpb.Marshaler
//	payload, err := m.MarshalToString(&reqParams)
//	assert.NoError(t, err, "failed to marshal to string")
//
//	// Without this running this on Travis CI might generate a connection refused error
//	// because the server may not be ready to accept connections just yet.
//	time.Sleep(3 * time.Second)
//
//	url := fmt.Sprintf("http://127.0.0.1:%d/v1/nonce", config.ConfigValues.JSONServerPort)
//	resp, err := http.Post(url, contentType, strings.NewReader(payload))
//	assert.NoError(t, err, "failed to http post to api endpoint")
//
//	defer resp.Body.Close()
//	_, err = ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	got, want := resp.StatusCode, http.StatusInternalServerError //todo: should we change it to err 400 somehow?
//	assert.Equal(t, want, got)
//
//	value := resp.Header.Get("Content-Type")
//	assert.Equal(t, value, contentType)
//
//	url = fmt.Sprintf("http://127.0.0.1:%d/v1/balance", config.ConfigValues.JSONServerPort)
//	resp, err = http.Post(url, contentType, strings.NewReader(payload))
//	assert.NoError(t, err, "failed to http post to api endpoint")
//
//	defer resp.Body.Close()
//
//	if got, want := resp.StatusCode, http.StatusInternalServerError; got != want {
//		t.Errorf("resp.StatusCode = %d; want %d", got, want)
//	}
//
//	value = resp.Header.Get("Content-Type")
//	assert.Equal(t, value, contentType)
//
//	// stop the services
//	jsonService.StopService()
//	<-jsonStatus
//	grpcService.StopService()
//	<-grpcStatus
//}
