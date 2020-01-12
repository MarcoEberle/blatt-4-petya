# Kommunikation zwischen den Services:
    - über Proto (gRPC) Nachrichten
    
![a](https://github.com/ob-vss-ws19/blatt-4-petya/blob/Develop/ProtoBild.PNG)

# Services
    - Erstellung eines Services
![a](https://github.com/ob-vss-ws19/blatt-4-petya/blob/Develop/NewServiceBild.PNG)

    - Kommunikation zwischen zwei Services
![b](https://github.com/ob-vss-ws19/blatt-4-petya/blob/Develop/KommunikationServicesBild.PNG)
# BookingService:
    - HallService
    - ShowService
    
# HallService:
    - ShowService
    
# MovieService:
    - ShowService
    
# ShowService:
    - MovieService
    - BookingService
    - HallService
    
# UserService:
    - BookingService
