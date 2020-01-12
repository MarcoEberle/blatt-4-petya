# Kommunikation zwischen den Services:
    - über Proto (gRPC) Nachrichten
    
    - Bsp.: ![](link)

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
