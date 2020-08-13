import React, { Component } from "react";
import aft from "../aft";
import Form from "@rjsf/core";
import { isObject, isFunction } from "../util";

class AftForm extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.handleSubmit = this.handleSubmit.bind(this);

    if (this.props.model !== "") {
      let load = aft.function.reactJsonSchemaForm({
        args: { model: this.props.model }
      });
      load.then(obj => {
        this.setState({ schema: obj["schema"], uiSchema: obj["uiSchema"] });
      });
    }
  }

  handleSubmit(e) {
    let errors = aft.function.validateForm({args: { schema: this.state.schema, data: e.formData }});
    errors.then(obj => {
		this.setState(prevState => {
			return {
			  ...prevState,
			  errors: obj,
			  formData: e.formData
			};
		});
		if(Object.keys(obj).length === 0) {
			let submit = aft.api[this.props.model].create({ data: e.formData });
			submit.then(obj => {
			    if(isFunction(this.props.handleSuccess)) {
	      			this.props.handleSuccess(obj);
	      		}
   			});
		} else {
			if (isFunction(this.props.handleError)) {
	      		this.props.handleError(obj);
	      	}
		}
  	});
  }
  render() {
  	var submit = "";
  	if (this.props.buttonClassName != null) {
  		var buttonText = this.props.buttonText != null? this.props.buttonText : "Submit";
  		submit = <button type="submit" className={this.props.buttonClassName}>{buttonText}</button>
  	}
    if (isObject(this.state.schema)) {
      return (
        <Form
          schema={this.state.schema}
          extraErrors={this.state.errors}
          formData={this.state.formData}
          uiSchema={this.state.uiSchema}
          onSubmit={fd => this.handleSubmit(fd)}
          ObjectFieldTemplate={ObjectFieldTemplate(this.props.fieldClassName)}
          FieldTemplate={
        	CustomFieldTemplate(
        		this.props.fieldTemplate, 
        		this.props.title, 
        		this.props.titleClassName,
        		this.props.description, 
        		this.props.descriptionClassName,
        		this.props.labelClassName,
        		this.props.inputClassName,
        	)
          }
          className={this.props.className}
        >{submit}</Form>
      );
    }
    return <div></div>;
  }
}

function ObjectFieldTemplate(cn) {
  return (props => {
    return (
      <div key={props.idSchema.$id + "wrap"}>
        {props.title}
        {props.description}
        {props.properties.map((element, idx) => <div key={props.idSchema.$id + "wrap" + idx} className={cn}>{element.content}</div>)}
      </div>
    );
  });
}

function CustomFieldTemplate(group, title, titleClassName, desc, descrClassName, labelClassName, inputClassName) {
  return (props => {
    const {id, classNames, label, help, required, errors, children} = props;
    if (!isFunction(group)) {
	  if (id === "root") {
		  return (<div key={"title"}>
		  			<div className={titleClassName} key={"title"}>{title}</div>
		  			<div className={descrClassName} key={ "desc"}>{desc}</div>
		  		  <div key={"inner"}>{children}</div>
		  		  </div>)      
	  }
      return (<div className={classNames} key={id + "outer"}>
        <label htmlFor={id} className={labelClassName}>{label}{required ? "*" : null}</label>
        <div className={inputClassName} key={id + "inner"}>
	        {children}
    	    {errors}
        	{help}
        </div>
      </div>)
    }
    return group(props);
  });
}

export default AftForm;
