// frontend/src/components/ThresholdSettings.jsx
import React, { useContext, useState, useEffect } from 'react';
import { 
    Button, 
    TextField, 
    Paper, 
    Typography, 
    Grid,
    Snackbar,
    Alert,
    CircularProgress
} from '@mui/material';
import { ThresholdsContext } from '../context/ThresholdsContext';

const TYPES = [
    { id: 'temperature', label: 'Температура', unit: '°C' },
    { id: 'humidity', label: 'Влажность', unit: '%' },
    { id: 'pressure', label: 'Давление', unit: 'hPa' }
];


export default function ThresholdSettings() {
    const { thresholds, updateThreshold,syncThresholds, error } = useContext(ThresholdsContext);
    const [localThresholds, setLocalThresholds] = useState({});
    const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });
    const [saving, setSaving] = useState({});

    useEffect(() => {
        if (thresholds) {
            setLocalThresholds(JSON.parse(JSON.stringify(thresholds)));
        }
    }, [thresholds]);

    const [isDirty, setIsDirty] = useState({});

    const handleChange = (type, field, value) => {
        const numValue = parseFloat(value);
        if (isNaN(numValue)) return;
        console.log(type, field, value);
        setIsDirty(prev => ({ ...prev, [type]: true }));

        setLocalThresholds(prev => ({
            ...prev,
            [type]: {
                ...prev[type],
                [field]: numValue
            }
        }));
    };

    useEffect(() => {
        if (thresholds && !Object.values(isDirty).some(Boolean)) {
            setLocalThresholds(JSON.parse(JSON.stringify(thresholds)));
        }
    }, [thresholds]);

    const handleSave = async (type) => {
        if (!localThresholds[type] || 
            localThresholds[type].min === undefined || 
            localThresholds[type].max === undefined) {
            setSnackbar({
                open: true,
                message: 'Пожалуйста, заполните оба значения',
                severity: 'error'
            });
            return;
        }
    
        setSaving(prev => ({ ...prev, [type]: true }));
        
        try {
            const response = await fetch('http://localhost:8080/api/thresholds/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ 
                    type, 
                    min_value: localThresholds[type].min, 
                    max_value: localThresholds[type].max 
                }),
            });
    
            if (!response.ok) {
                throw new Error(await response.text());
            }

            const updatedResponse = await fetch('http://localhost:8080/api/thresholds');
            const updatedData = await updatedResponse.json();
            
            const thresholdsMap = {};
            updatedData.forEach(t => {
                thresholdsMap[t.type] = { 
                    min: t.min_value, 
                    max: t.max_value 
                };
            });
            
            // Обновляем состояние через syncThresholds
            syncThresholds(thresholdsMap);
    
            setSnackbar({
                open: true,
                message: 'Пороги успешно обновлены',
                severity: 'success'
            });
        } catch (error) {
            setSnackbar({
                open: true,
                message: `Ошибка: ${error.message}`,
                severity: 'error'
            });
        } finally {
            setSaving(prev => ({ ...prev, [type]: false }));
        }
    };

    const handleCloseSnackbar = () => {
        setSnackbar(prev => ({ ...prev, open: false }));
    };

    if (!thresholds || !localThresholds) return <CircularProgress />;

    return (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
                Настройка пороговых значений
            </Typography>
            
            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    Ошибка загрузки порогов: {error}. Используются значения по умолчанию.
                </Alert>
            )}
            
            <Grid container spacing={3}>
                {TYPES.map(({ id, label, unit }) => (
                    <Grid item xs={12} md={4} key={id}>
                        <Typography variant="subtitle1" gutterBottom>
                            {label} ({unit})
                        </Typography>
                        
                        <TextField
                            label="Минимальное значение"
                            type="number"
                            value={localThresholds[id]?.min || ''}
                            onChange={(e) => handleChange(id, 'min', e.target.value)}
                            fullWidth
                            margin="normal"
                            InputLabelProps={{ shrink: true }}
                        />
                        
                        <TextField
                            label="Максимальное значение"
                            type="number"
                            value={localThresholds[id]?.max || ''}
                            onChange={(e) => handleChange(id, 'max', e.target.value)}
                            fullWidth
                            margin="normal"
                            InputLabelProps={{ shrink: true }}
                        />
                        
                        <Button
                            variant="contained"
                            onClick={() => handleSave(id)}
                            disabled={saving[id]}
                            startIcon={saving[id] ? <CircularProgress size={20} /> : null}
                            sx={{ mt: 1 }}
                        >
                            {saving[id] ? 'Сохранение...' : 'Сохранить'}
                        </Button>
                    </Grid>
                ))}
            </Grid>
            
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={handleCloseSnackbar}
            >
                <Alert 
                    onClose={handleCloseSnackbar} 
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Paper>
    );
}