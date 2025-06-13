import grpc from "k6/net/grpc";
import { check } from "k6";

const client = new grpc.Client();
client.load(["../../proto"], "auth.proto");

export let options = {
  vus: 1000,
  duration: "10s",
};

export default () => {
  client.connect("localhost:50051", {
    plaintext: true,
  });

  const response = client.invoke("auth.AuthService/Login", {
    username: "testuser",
    password: "testpass",
  });

  check(response, {
    "status is OK": (r) => r && r.status === grpc.StatusOK,
  });

  client.close();
};
