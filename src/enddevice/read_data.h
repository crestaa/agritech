enum humidity {
  PIN = 1,
  MAX = 4000,
  MIN = 2000,
  MIN_LIMIT = 1500,
  MAX_LIMIT = 4500
};

int getHumidity() {
  int value = analogRead(humidity::PIN);

  // disconnected or malfunctioning sensor
  if (value < humidity::MIN_LIMIT || value > humidity::MAX_LIMIT) return -1;
  
  // safe zone for 0% humidity
  if (value >= humidity::MAX) return 0;
  // safe zone for 100% humidity
  else if (value <= humidity::MIN) return 100; 
  // mapping value MIN->100%, MAX->0%
  else return (int)map(value, humidity::MAX, humidity::MIN, 0, 100);
}
