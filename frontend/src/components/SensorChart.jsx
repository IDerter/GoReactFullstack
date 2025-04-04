import React, { useContext, useEffect, useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import { 
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler
} from 'chart.js';
import Annotation from 'chartjs-plugin-annotation';
import { ThresholdsContext } from '../context/ThresholdsContext';
import PropTypes from 'prop-types';

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler,
    Annotation
);

export default function SensorChart({ data, type }) {
    const { thresholds } = useContext(ThresholdsContext);

    useEffect(() => {
        console.log(`Thresholds updated for ${type}:`, thresholds?.[type]);
    }, [thresholds, type]);

    // Получаем текущие пороги с проверкой
    const currentThresholds = thresholds?.[type] || {};
    const safeMin = Math.min(
        typeof currentThresholds.min === 'number' ? currentThresholds.min : 0,
        typeof currentThresholds.max === 'number' ? currentThresholds.max : 100
    );
    const safeMax = Math.max(
        typeof currentThresholds.min === 'number' ? currentThresholds.min : 0,
        typeof currentThresholds.max === 'number' ? currentThresholds.max : 100
    );

    // Вычисляем min/max для шкалы
    const allDataPoints = data?.datasets?.flatMap(d => d.data || []) || [];
    const minValue = allDataPoints.length > 0 ? 
        Math.min(safeMin, ...allDataPoints) - 5 : safeMin - 5;
    const maxValue = allDataPoints.length > 0 ? 
        Math.max(safeMax, ...allDataPoints) + 5 : safeMax + 5;

    // Опции графика с использованием useMemo
    const options = useMemo(() => ({
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'top',
            },
            tooltip: {
                callbacks: {
                    label: (context) => {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y;
                        let warning = '';
                        
                        if (value < safeMin) warning = ' (ниже минимума)';
                        if (value > safeMax) warning = ' (выше максимума)';
                        
                        return `${label}: ${value} ${warning}`;
                    }
                }
            },
            annotation: {
                annotations: {
                    minLine: {
                        type: 'line',
                        yMin: safeMin,
                        yMax: safeMin,
                        borderColor: 'rgba(255, 99, 132, 0.7)',
                        borderWidth: 2,
                        borderDash: [6, 6],
                        label: {
                            content: `Min: ${safeMin.toFixed(1)}`,
                            enabled: true,
                            position: 'left',
                            backgroundColor: 'rgba(255, 99, 132, 0.5)'
                        }
                    },
                    maxLine: {
                        type: 'line',
                        yMin: safeMax,
                        yMax: safeMax,
                        borderColor: 'rgba(54, 162, 235, 0.7)',
                        borderWidth: 2,
                        borderDash: [6, 6],
                        label: {
                            content: `Max: ${safeMax.toFixed(1)}`,
                            enabled: true,
                            position: 'left',
                            backgroundColor: 'rgba(54, 162, 235, 0.5)'
                        }
                    },
                    dangerZoneMin: {
                        type: 'box',
                        yMin: -Infinity,
                        yMax: safeMin,
                        backgroundColor: 'rgba(255, 0, 0, 0.1)',
                        borderWidth: 0
                    },
                    dangerZoneMax: {
                        type: 'box',
                        yMin: safeMax,
                        yMax: Infinity,
                        backgroundColor: 'rgba(255, 0, 0, 0.1)',
                        borderWidth: 0
                    }
                }
            }
        },
        scales: {
            y: {
                min: minValue,
                max: maxValue,
                ticks: {
                    callback: (value) => {
                        if (value < safeMin || value > safeMax) {
                            return `${value}*`;
                        }
                        return value;
                    },
                    color: (context) => {
                        const value = context.tick.value;
                        if (value < safeMin) return 'red';
                        if (value > safeMax) return 'red';
                        return 'black';
                    }
                }
            }
        }
    }), [safeMin, safeMax, minValue, maxValue]);

    // Проверка данных после всех хуков
    if (!data || !Array.isArray(data.labels) || !data.datasets?.length) {
        return (
            <div style={{ 
                height: '300px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                border: '1px dashed #ccc',
                color: '#666'
            }}>
                Ожидание данных...
            </div>
        );
    }

    if (!thresholds?.[type]) {
        return <div>Загрузка пороговых значений для {type}...</div>;
    }
    
    return (
        <div style={{ height: '400px', position: 'relative' }}>
            <Line 
                key={`chart-${type}-${thresholds?.[type]?.min}-${thresholds?.[type]?.max}`}
                data={data} 
                options={options}
                redraw
            />
        </div>
    );
}

SensorChart.propTypes = {
    data: PropTypes.shape({
        labels: PropTypes.array.isRequired,
        datasets: PropTypes.arrayOf(
            PropTypes.shape({
                label: PropTypes.string,
                data: PropTypes.array.isRequired
            })
        ).isRequired
    }),
    type: PropTypes.oneOf(['temperature', 'humidity', 'pressure']).isRequired
};