import React from 'react';
import Select from 'react-select';
import CreatableSelect from 'react-select/creatable';
import classNames from 'classnames';

export function MultiSelect({options, className, onChange, defaultValue}) {

   return (
        <Select
            defaultValue={defaultValue}
            onChange={onChange}
            isMulti
            name="def"
            options={options}
            className={classNames(className, "basic-multi-select")}
            classNamePrefix="select"
        />
    )
}

export function SelectCreatable({options, className, onChange}) {

    return (
         <CreatableSelect
             //defaultValue={}
             onChange={onChange}
             isMulti
             isClearable
             name="def"
             options={options}
             className={classNames(className, "basic-multi-select")}
             classNamePrefix="select"
         />
     )
 }
