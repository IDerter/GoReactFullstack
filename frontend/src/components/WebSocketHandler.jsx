import React, { useEffect, useContext } from 'react';
import { ThresholdsContext } from '../context/ThresholdsContext';

export default function WebSocketHandler({ onDataReceived }) {
    const { syncThresholds } = useContext(ThresholdsContext);

    useEffect(() => {
        const ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        ws.onopen = () => {
            console.log('WebSocket connected');
        };
        
        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                console.log('Raw WebSocket message:', message, onDataReceived);
               
                // Обработка sensor data
                if (message.data && onDataReceived) {
                    console.log('Processing sensor data...');
                    const chartData = processSensorData(message.data);
                    console.log('Processed chart data:', chartData);
                    onDataReceived(chartData);
                }
            } catch (error) {
                console.error('Error processing WebSocket message:', error);
            }
        };
        
        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        ws.onclose = () => {
            console.log('WebSocket disconnected');
        };
        
        return () => {
            ws.close();
        };
    }, [syncThresholds, onDataReceived]);

    return null;
}

function processSensorData(sensorData) {
    const result = {};
    
    sensorData.forEach(item => {
        if (!result[item.type]) {
            result[item.type] = {
                labels: [],
                datasets: [{
                    label: item.type,
                    data: [],
                    borderColor: getColorForType(item.type),
                    backgroundColor: 'rgba(0, 0, 0, 0.1)'
                }]
            };
        }
        result[item.type].labels.push(new Date(item.timestamp).toLocaleTimeString());
        result[item.type].datasets[0].data.push(item.value);
    });
    
    return result;
}

function getColorForType(type) {
    switch(type) {
        case 'temperature': return 'rgb(255, 99, 132)';
        case 'humidity': return 'rgb(54, 162, 235)';
        case 'pressure': return 'rgb(75, 192, 192)';
        default: return 'rgb(201, 203, 207)';
    }
}