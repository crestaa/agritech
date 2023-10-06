const field_id = 1;
var field;
var sensors;
var readings;

function fetchField(){
    var url = '../api/fields/'+field_id
    fetch(url,{}).then((response)=>response.json()).then(
        (data)=>{
            console.log(data)
            field = data

            document.getElementById('field_name').innerHTML+=" "+field.Nome
            var coords = []
            var c = { lon: field.Lon, lat: field.Lat }
            coords.push(c) 
            initMap('map',coords)
        }
    )
}
fetchField()

function fetchSensors(){
    var url = '../api/fields/'+field_id+"/sensors"
    fetch(url,{}).then((response)=>response.json()).then(
        (data)=>{
            console.log(data)
            sensors = data

            sensors.forEach(function(el){
                document.getElementById('sensorsList').innerHTML += [{ id:el.ID, mac:el.MAC, lat:el.Lat, lon:el.Lon }].map(sensorItem)
            })
        }
    )
}
fetchSensors()

function fetchReadings(){
    var url = '../api/fields/'+field_id+"/readings"
    fetch(url,{}).then((response)=>response.json()).then(
        (data)=>{
            console.log(data)
            readings = data

            readings.forEach(function(el){
                v=el.Valore
                t=el.ID_tipo_misurazione
                if(t==1){t="Humidity"; v+="%"}
                document.getElementById('readingsList').innerHTML += [{ sens_id:el.ID_sensore, type:t, value:v, time:el.Data_ora }].map(readingItem)
            })
        }
    )
}
fetchReadings()

//item template
const sensorItem = ({ id, mac, lat, lon }) => `
    <tr p-2 class="py-8 border-b border-solid border-gray-300">
        <td class="p-2 py-4 border-b border-solid border-gray-300">
            <div class="pl-4 flex flex-wrap flex-row items-center">
                <div class="mr-4 h-16 w-32 block flex flex-row items-center">${id}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${mac}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${lat}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${lon}</div>
            </div>
        </td>
    </tr>
`
const readingItem = ({ sens_id, type, value, time }) => `
    <tr p-2 class="py-8 border-b border-solid border-gray-300">
        <td class="p-2 py-4 border-b border-solid border-gray-300">
            <div class="pl-4 flex flex-wrap flex-row items-center">
                <div class="mr-4 h-16 w-32 block flex flex-row items-center">${sens_id}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${type}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${value}</div>
                <div class="mr-4 h-16 w-64 block flex flex-row items-center">${time}</div>
            </div>
        </td>
    </tr>
`

// MAP-RELATED FUNCTIONS
function initMap(target, coords) {
  map = new OpenLayers.Map(target)
  map.addLayer(new OpenLayers.Layer.OSM())

  var center = new OpenLayers.LonLat(11.3426163, 44.494887).transform(
    new OpenLayers.Projection('EPSG:4326'),
    map.getProjectionObject()
  )

  var zoom = 8

  map.setCenter(center, zoom)

  addMarkers(map, coords)
}
function addMarkers(map, coords) {
  var markers = new OpenLayers.Layer.Markers('Markers')
  var icon = new OpenLayers.Icon(
    'https://icons.iconarchive.com/icons/paomedia/small-n-flat/32/map-marker-icon.png',
    new OpenLayers.Size(25, 25),
    new OpenLayers.Pixel(-(25 / 2), -25)
  )

  map.addLayer(markers)
  var locations = new Array()
  coords.forEach((c) => {
    locations.push(
      new OpenLayers.LonLat(c.lon, c.lat).transform(new OpenLayers.Projection('EPSG:4326'), map.getProjectionObject())
    )
  })

  locations.forEach((loc) => {
    markers.addMarker(new OpenLayers.Marker(loc, icon.clone()))
  })
}