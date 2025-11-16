import React, { Suspense } from 'react';
import Typography from '@mui/material/Typography';

const OAuth2 = React.lazy(() => import('./fields/oauth2/OAuth2'));
const PhotoSelect = React.lazy(() => import('./fields/photoselect/PhotoSelect'));
const RawPhotoSelect = React.lazy(() => import('./fields/photoselect/RawPhotoSelect'));
const Toggle = React.lazy(() => import('./fields/Toggle'));
const Color = React.lazy(() => import('./fields/Color'));
const DateTime = React.lazy(() => import('./fields/DateTime'));
const Dropdown = React.lazy(() => import('./fields/Dropdown'));
const LocationBased = React.lazy(() => import('./fields/location/LocationBased'));
const LocationForm = React.lazy(() => import('./fields/location/LocationForm'));
const TextInput = React.lazy(() => import('./fields/TextInput'));
const Typeahead = React.lazy(() => import('./fields/Typeahead'));

export default function FieldDetails({ field }) {
    return (
        <Suspense fallback={<div />}> 
            {(() => {
                switch (field.type) {
                    case 'datetime':
                        return <DateTime field={field} />
                    case 'dropdown':
                        return <Dropdown field={field} />
                    case 'location':
                        return <LocationForm field={field} />
                    case 'locationbased':
                        return <LocationBased field={field} />
                    case 'oauth2':
                        return <OAuth2 field={field} />
                    case 'png':
                        return <PhotoSelect field={field} />
                    case 'text':
                        return <TextInput field={field} />
                    case 'onoff':
                        return <Toggle field={field} />
                    case 'typeahead':
                        return <Typeahead field={field} />
                    case 'color':
                        return <Color field={field} />
                    default:
                        return <Typography>Unsupported type: {field.type}</Typography>
                }
            })()}
        </Suspense>
    )
}