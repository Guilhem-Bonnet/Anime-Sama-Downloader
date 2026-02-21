import React, { useState } from 'react';

export interface FormStepperStep {
  id: string;
  label: string;
  completed: boolean;
}

export interface FormStepperProps {
  steps: FormStepperStep[];
  currentStep: number;
  onStepClick?: (stepIndex: number) => void;
}

export const FormStepper: React.FC<FormStepperProps> = ({ steps, currentStep, onStepClick }) => {
  return (
    <div className="w-full">
      <div className="flex items-center justify-between mb-8">
        {steps.map((step, index) => (
          <React.Fragment key={step.id}>
            <button
              onClick={() => onStepClick?.(index)}
              className={`flex flex-col items-center ${
                onStepClick ? 'cursor-pointer' : 'cursor-default'
              }`}
            >
              <div
                className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold transition-colors ${
                  index === currentStep
                    ? 'bg-cyan-500 text-white'
                    : index < currentStep
                      ? 'bg-green-500 text-white'
                      : 'bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300'
                }`}
              >
                {index < currentStep ? '✓' : index + 1}
              </div>
              <span
                className={`text-xs mt-2 text-center text-gray-600 dark:text-gray-400 ${
                  index === currentStep ? 'font-semibold text-cyan-600 dark:text-cyan-400' : ''
                }`}
              >
                {step.label}
              </span>
            </button>

            {index < steps.length - 1 && (
              <div
                className={`flex-1 h-1 mx-2 ${
                  index < currentStep - 1
                    ? 'bg-green-500'
                    : index < currentStep
                      ? 'bg-cyan-500'
                      : 'bg-gray-300 dark:bg-gray-600'
                }`}
              />
            )}
          </React.Fragment>
        ))}
      </div>
    </div>
  );
};
