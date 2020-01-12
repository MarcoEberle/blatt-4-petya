Anleitung:

1. Starten Sie alle Servies. Dazu können Sie die StartAllServices.sh benutzen.

2. Starten Sie client.go

Anmerkung zum Client:

Der Client führt das verlangte Szenario komplett aus.
Folgende Schritte werden ausgeführt:
    1. Erstellung von 4 Filmen
    2. Erstellung von 2 Sälen
    3. Erstellung von 4 Benutzern
    4. Erstellung von 4 Vorstellungen
        (Jeder Film wird in einer Vorstellung gezeigt)
    5. Erstellung von 4 Buchungen
        Dabei werden die ersten zwei Buchungen gleichzeitig ausgeführt.
        Nur eine Buchung wird erstellt.
        Desweiteren werden 3 weitere Buchungen erstellt.
    6. Bestätigung von 3 Buchungen
    7. Löschung eines Saals
        Im BookingService werden die gespeicherten Bookings ausgegeben.
        Dort sehen Sie, dass von den vier Bookings nur noch zwei übrig bleiben.
        ShowID1 und ShowID2 sind in Halle HallID1
        Alle Shows und somit Bookings der HallID2 wurden gelöscht.
    
Alle Prozesse werden im Client geloggt bzw. ausführlich in den jeweiligen Services geloggt.
