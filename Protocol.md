Kommunikation zwischen den Services:
    - über Proto (gRPC) Nachrichten

BookingService:
    - HallService
    - ShowService
    
HallService:
    - ShowService
    
MovieService:
    - ShowService
    
ShowService:
    - MovieService
    - BookingService
    - HallService
    
UserService:
    - BookingService