#!/bin/bash

# set date and time
date_time=$(date +"%Y_%m_%d_%H%M%S")

# backup previous weights
mv ../email-processing-webapp/weights/gemma2-2b_inquiry_tuned.lora.h5 ../email-processing-webapp/weights/gemma2-2b_inquiry_tuned.lora.h5.$date_time.backup

# deploy new weights
cp weights/gemma2-2b_inquiry_tuned_4_epoch3.lora.h5 ../email-processing-webapp/weights/gemma2-2b_inquiry_tuned.lora.h5