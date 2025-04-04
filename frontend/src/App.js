// frontend/src/App.js
import React, { useState } from 'react';
import { ThresholdsProvider } from './context/ThresholdsContext';
import ThresholdSettings from './components/ThresholdSettings';
import SensorChart from './components/SensorChart';
import WebSocketHandler from './components/WebSocketHandler';

function App() {
    const [chartData, setChartData] = useState({
        temperature: emptyChartData('temperature'),
        humidity: emptyChartData('humidity'),
        pressure: emptyChartData('pressure')
    });

    const handleDataReceived = (newData) => {
        console.log('Received new data:', newData);
        setChartData(prev => ({
            ...prev,
            ...newData
        }));
    };

    return (
        <ThresholdsProvider>
            <div style={{ padding: '20px' }}>
                <WebSocketHandler onDataReceived={handleDataReceived} />
                
                <ThresholdSettings />
                
                <div style={{ marginTop: '30px' }}>
                    <h2>Графики показаний</h2>
                    
                    {chartData.temperature && (
                        <div style={{ marginBottom: '30px' }}>
                            <h3>Температура</h3>
                            <SensorChart data={chartData.temperature || emptyChartData('temperature')} type="temperature" />
                        </div>
                    )}
                    
                    {chartData.humidity && (
                        <div style={{ marginBottom: '30px' }}>
                            <h3>Влажность</h3>
                            <SensorChart data={chartData.humidity || emptyChartData('humidity')} type="humidity" />
                        </div>
                    )}
                    
                    {chartData.pressure && (
                        <div style={{ marginBottom: '30px' }}>
                            <h3>Давление</h3>
                            <SensorChart data={chartData.pressure || emptyChartData('pressure')} type="pressure" />
                        </div>
                    )}
                </div>
            </div>
            
        </ThresholdsProvider>
    );
    function emptyChartData(type) {
        return {
            labels: [],
            datasets: [{
                label: type,
                data: [],
                borderColor: getColorForType(type),
                backgroundColor: 'rgba(0, 0, 0, 0.1)'
            }]
        };
    }
    function getColorForType(type) {
        switch(type) {
            case 'temperature': return 'rgb(255, 99, 132)';
            case 'humidity': return 'rgb(54, 162, 235)';
            case 'pressure': return 'rgb(75, 192, 192)';
            default: return 'rgb(201, 203, 207)';
        }
    }
    
}

export default App;