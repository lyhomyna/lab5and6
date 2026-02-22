#!/bin/bash

APP_URL="http://localhost:8080/parts"

echo "Starting to populate the JDM Registry with 11 legendary items..."

parts=("RB26DETT Engine" "2JZ-GTE Engine" "TE37 Wheels" "Brembo Brakes" "Momo Steering Wheel" "Recaro Seats" "Tomei Expreme Exhaust" "HKS Turbo Kit" "Ohlins Suspension" "Nismo Body Kit" "Greddy Intercooler")
cars=("Nissan Skyline" "Toyota Supra" "Nissan Silvia" "Mitsubishi Evo" "Honda NSX" "Mazda RX-7" "Subaru Impreza" "Toyota AE86" "Nissan 350Z" "Honda Civic Type R" "Nissan GT-R")

for i in {0..10}
do
    name=$(echo ${parts[$i]} | sed 's/ /%20/g')
    model=$(echo ${cars[$i]} | sed 's/ /%20/g')
    
    echo "Adding item $((i+1)): ${parts[$i]} for ${cars[$i]}"
    
    curl -s -X POST "$APP_URL?name=$name&model=$model"
    echo -e "\n"
done

echo "Done! 11 items added to the cloud database."
