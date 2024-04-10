const map = L.map('map').setView([0, 0], 2);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors',
  maxZoom: 18,
}).addTo(map);

async function fetchUserData(userId) {
  const response = await fetch(`/${userId}/data`);
  const data = await response.json();
  return data;
}

function getCountryColor(activityNumber, sumDistance) {
  if (activityNumber > 0) {
    return 'green';
  }
  return 'gray';
}

function formatActivityInfo(activityNumber, sumDistance) {
  const distanceInKm = (sumDistance / 1000).toFixed(0);
  return `Number of activities: ${activityNumber}<br>Total distance: ${distanceInKm}km`;
}

const userId = window.location.pathname.split('/')[1];

fetchUserData(userId)
  .then(data => {
    fetch('/static/ne_50m_admin_0_countries.json')
      .then(response => response.json())
      .then(geojson => {
        const countriesLayer = L.geoJSON(geojson, {
          style: feature => {
            const countryCode = feature.properties.ADM0_A3;
            const countryData = data['countries'][countryCode];
            if (countryData) {
              const { ActivityNumber, SumDistance } = countryData;
              return {
                fillColor: getCountryColor(ActivityNumber, SumDistance),
                weight: 1,
                opacity: 1,
                color: 'white',
                fillOpacity: 0.7,
              };
            }
            return {
              fillColor: 'gray',
              weight: 1,
              opacity: 1,
              color: 'white',
              fillOpacity: 0.7,
            };
          },
          onEachFeature: (feature, layer) => {
            const countryCode = feature.properties.ADM0_A3;
            const countryData = data['countries'][countryCode];
            if (countryData) {
              const { ActivityNumber, SumDistance } = countryData;
              layer.bindPopup(formatActivityInfo(ActivityNumber, SumDistance));
            }
          },
        }).addTo(map);
      });
  });
