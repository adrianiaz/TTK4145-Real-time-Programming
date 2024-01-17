#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define PORT 30000
#define BUFFER_SIZE 80

//This file recieves the value from the UDP server of the local network, it should display the message "Hello from UDP server at 10.100.23.129!"


int main() {
    int socketDescriptor;
    char storageBuffer[BUFFER_SIZE];
    int lengthOfAddress;
    struct sockaddr_in recieverAddress;
    struct sockaddr_in connectionAddress;

    // Creating a socket with IPv4 domain and UDP protocol
    socketDescriptor = socket(AF_INET, SOCK_DGRAM, 0);
    if (socketDescriptor < 0) {
        perror("Socket creation failed");
        exit(EXIT_FAILURE);
    }

    // Set options for the socket. SO_REUSEADDR allows local addresses to be reused in the bind() function
    int optionValue = 1;
    if (setsockopt(socketDescriptor, SOL_SOCKET, SO_REUSEADDR, &optionValue, sizeof(optionValue)) < 0) {
        perror("Couldn't set socket options");
        exit(EXIT_FAILURE);
    }

    // Initialize structure elements for address
    recieverAddress.sin_family = AF_INET; // IPv4
    recieverAddress.sin_port = htons(PORT);
    recieverAddress.sin_addr.s_addr = INADDR_ANY;
    memset(recieverAddress.sin_zero, '\0', sizeof(recieverAddress.sin_zero));

    // Bind the socket
    if (bind(socketDescriptor, (struct sockaddr*)&recieverAddress, sizeof(struct sockaddr)) < 0) {
        perror("Couldn't bind socket");
        exit(EXIT_FAILURE);
    }

    // Receive data sent by lab server
    lengthOfAddress = sizeof(connectionAddress);
    ssize_t bytesReceived = recvfrom(socketDescriptor, storageBuffer, sizeof(storageBuffer), 0, (struct sockaddr*)&connectionAddress, &lengthOfAddress);
    if (bytesReceived < 0) {
        perror("Error receiving data");
        exit(EXIT_FAILURE);
    }

    // Set the last index of the character array as a null character
    storageBuffer[bytesReceived] = '\0';
    printf("\n ---incoming message from labserver---\n\n %s \n", storageBuffer);
    printf("\n-----------------------------------------\n\n");

    // Close the socket
    close(socketDescriptor);

    return 0;
}
