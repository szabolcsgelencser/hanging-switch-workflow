#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

#include <netinet/tcp.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <netdb.h>

#include <pthread.h>

void *fork_and_exec(void *vargp)
{
	pid_t pid = fork();
	if (pid == -1) {} // failed to fork 
	else if (pid > 0)
	{
        // we are the parent
		int status;
		waitpid(pid, &status, 0);
	}
	else 
	{
		// we are the child
		execve("./forkexec.exe", NULL, NULL);
		_exit(EXIT_FAILURE);   // exec never returns
	}
    return NULL;
}

int socket_connect(char *host, in_port_t port){
	struct hostent *hp;
	struct sockaddr_in addr;
	int on = 1, sock;     

	if((hp = gethostbyname(host)) == NULL){
		herror("gethostbyname");
		exit(1);
	}
	bcopy(hp->h_addr, &addr.sin_addr, hp->h_length);
	addr.sin_port = htons(port);
	addr.sin_family = AF_INET;
	sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
	setsockopt(sock, IPPROTO_TCP, TCP_NODELAY, (const char *)&on, sizeof(int));

	if(sock == -1){
		perror("setsockopt");
		exit(1);
	}
	
	if(connect(sock, (struct sockaddr *)&addr, sizeof(struct sockaddr_in)) == -1){
		perror("connect");
		exit(1);

	}
	return sock;
}
 
#define BUFFER_SIZE 1024

void *http_request(void *vargp)
{
    int fd;
	char buffer[BUFFER_SIZE];

	fd = socket_connect("google.com", 80);
	write(fd, "GET /\r\n", strlen("GET /\r\n")); // write(fd, char[]*, len);
	bzero(buffer, BUFFER_SIZE);
	
	while(read(fd, buffer, BUFFER_SIZE - 1) != 0){
		// fprintf(stderr, "%s", buffer);
		bzero(buffer, BUFFER_SIZE);
	}

	shutdown(fd, SHUT_RDWR);
	close(fd);

    return NULL;
}

int main()
{
    pthread_t fork_exec_threads[100];
    for (int i = 0; i < 100; i++) {
        pthread_t thread_id;
        pthread_create(&thread_id, NULL, fork_and_exec, NULL);
        fork_exec_threads[i] = thread_id;
    }

    pthread_t http_threads[20];
    for (int i = 0; i < 20; i++) {
        pthread_t thread_id;
        pthread_create(&thread_id, NULL, http_request, NULL);
        http_threads[i] = thread_id;
    }
    
    for (int i = 0; i < 100; i++) {
        pthread_join(fork_exec_threads[i], NULL);
    }
    for (int i = 0; i < 20; i++) {
        pthread_join(http_threads[i], NULL);
    }
}
