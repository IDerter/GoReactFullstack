// frontend/src/context/ThresholdsContext.jsx
import React, { createContext, useState, useEffect } from 'react';

export const ThresholdsContext = createContext();

export function ThresholdsProvider({ children }) {
    const [thresholds, setThresholds] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);


    const syncThresholds = (newThresholds) => {
        // Просто обновляем состояние без отправки на сервер
        setThresholds(newThresholds);
    };

    const fetchThresholds = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/thresholds');
            if (!response.ok) throw new Error('Failed to fetch thresholds');
            const data = await response.json();
            
            const thresholdsMap = {};
            data.forEach(t => {
                thresholdsMap[t.type] = { 
                    min: t.min_value, 
                    max: t.max_value 
                };
            });
            
            setThresholds(thresholdsMap);
        } catch (err) {
            setError(err.message);
            // Установите дефолтные значения при ошибке
            setThresholds({
                temperature: { min: 20, max: 100 },
                humidity: { min: 30, max: 80 },
                pressure: { min: 900, max: 1100 }
            });
        } finally {
            setLoading(false);
        }
    };

    const updateThreshold = async (type, newValues) => {
        try {
            console.log(type, newValues)
            // Валидация
    
            if (newValues.min === undefined || newValues.max === undefined) {
                throw new Error('Не указаны min или max значения');
            }
    
            if (typeof newValues.min !== 'number' || typeof newValues.max !== 'number') {
                throw new TypeError('min и max должны быть числами');
            }
    
            if (newValues.min >= newValues.max) {
                throw new Error('min_value должен быть меньше max_value');
            }
    
            // Отправка запроса
            const response = await fetch('http://localhost:8080/api/thresholds/update', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({
                    type,
                    min_value: newValues.min,
                    max_value: newValues.max
                })
            });
    
            // Проверка ответа
            const contentType = response.headers.get('content-type');
            if (!contentType || !contentType.includes('application/json')) {
                const text = await response.text();
                throw new Error(`Ожидался JSON, но получен: ${contentType}. Ответ: ${text.substring(0, 100)}`);
            }
    
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка при обновлении порогов');
            }
    
            // Обновление состояния
            setThresholds(prev => ({
                ...prev,
                [type]: {
                    min: Number(newValues.min), 
                    max: Number(newValues.max) 
                }
            }));
    
            return { success: true, data };
        } catch (err) {
            console.error('Update threshold error:', err);
            return { 
                success: false, 
                error: err.message,
                stack: process.env.NODE_ENV === 'development' ? err.stack : undefined
            };
        }
    };

    useEffect(() => {
        fetchThresholds();
    }, []);

    if (loading) return <div>Loading thresholds...</div>;
    if (error) console.error('Thresholds loading error:', error);

    return (
        <ThresholdsContext.Provider value={{ thresholds, updateThreshold, syncThresholds, error }}>
            {children}
        </ThresholdsContext.Provider>
    );
}