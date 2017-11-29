#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>
#include <netinet/tcp.h>
#include <sys/time.h>

#define BUFFER_SIZE 1
#define on_error(...) { fprintf(stderr, __VA_ARGS__); fflush(stderr); exit(1); }
#define log_message(...) { fprintf(stderr, __VA_ARGS__); fflush(stderr); }

int main (int argc, char *argv[]) {
  if (argc < 2) on_error("Usage: %s [port]\n", argv[0]);

  int port = atoi(argv[1]);

  int server_fd, client_fd, err;
  struct sockaddr_in server, client;
  char buf[BUFFER_SIZE];

  server_fd = socket(AF_INET, SOCK_STREAM, 0);
  if (server_fd < 0) on_error("Could not create socket\n");

  server.sin_family = AF_INET;
  server.sin_port = htons(port);
  server.sin_addr.s_addr = htonl(INADDR_ANY);

  int opt_val = 1;
  setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt_val, sizeof opt_val);

  err = bind(server_fd, (struct sockaddr *) &server, sizeof(server));
  if (err < 0) on_error("Could not bind socket\n");

  err = listen(server_fd, 128);
  if (err < 0) on_error("Could not listen on socket\n");

  printf("Server is listening on %d\n", port);

  int flag = 1;
  while (1) {
    socklen_t client_len = sizeof(client);
    client_fd = accept(server_fd, (struct sockaddr *) &client, &client_len);
    /*if (setsockopt(client_fd, IPPROTO_TCP, TCP_NODELAY, &flag, sizeof(int)) == -1) {
        on_error("set TCP_NODELAY fail");
    }*/
    /*if (setsockopt(client_fd, IPPROTO_TCP, TCP_CORK, &flag, sizeof(int)) == -1) {
        on_error("set TCP_CORK fail");
    }*/

    if (client_fd < 0) on_error("Could not establish new connection\n");

    int i = 0;
    while (1) {
      i++;
      buf[0] = 'a';
      struct timeval stop, start;
      gettimeofday(&start, NULL);
      err = send(client_fd, buf, 1, 0);
      if (err < 0) on_error("Client write failed\n");
      gettimeofday(&stop, NULL);
      printf("[%d]took %lu\n", i, stop.tv_usec - start.tv_usec);
      //sleep(1);
    }
  }

  return 0;
}