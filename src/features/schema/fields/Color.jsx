import React, { useState, useEffect } from 'react';

import { useSelector, useDispatch } from 'react-redux';
import { HexColorPicker } from 'react-colorful';
import { set } from '../../config/configSlice';

export default function Color({ field }) {
    const initial = (field.default && field.default.startsWith('#')) ? field.default : (field.default ? `#${field.default}` : '#ffffff');
    const [color, setColor] = useState(initial);
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            const val = config[field.id].value;
            setColor(val && val.startsWith('#') ? val : (val ? `#${val}` : '#ffffff'));
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [config])

    const onChange = (value) => {
        setColor(value);
        dispatch(set({
            id: field.id,
            value: value,
        }));
    }

    return (
        <div>
            <HexColorPicker color={color} onChange={onChange} />
            {Array.isArray(field.palette) && field.palette.length > 0 && (
                <div style={{ display: 'flex', flexWrap: 'wrap', marginTop: 8 }}>
                    {field.palette.map((p, i) => {
                        const hex = (typeof p === 'string') ? (p.startsWith('#') ? p : `#${p}`) : (p && p.hex ? (p.hex.startsWith('#') ? p.hex : `#${p.hex}`) : '#000000');
                        const selected = hex.toLowerCase() === (color || '').toLowerCase();
                        return (
                            <button
                                key={i}
                                onClick={() => onChange(hex)}
                                title={hex}
                                style={{
                                    width: 28,
                                    height: 28,
                                    marginRight: 6,
                                    marginBottom: 6,
                                    borderRadius: 4,
                                    border: selected ? '2px solid #000' : '1px solid #ccc',
                                    background: hex,
                                    padding: 0,
                                    cursor: 'pointer'
                                }}
                                aria-pressed={selected}
                            />
                        )
                    })}
                </div>
            )}
        </div>
    )
}