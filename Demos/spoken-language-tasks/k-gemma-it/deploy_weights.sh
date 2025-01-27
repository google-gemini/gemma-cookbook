#!/bin/bash

# set date and time
date_time=$(date +"%Y_%m_%d_%H%M%S")

# backup previous weights
mv ../k-mail-replier/k_mail_replier/weights/gemma2-2b_k-tuned.lora.h5 ../k-mail-replier/k_mail_replier/weights/gemma2-2b_k-tuned.lora.h5.$date_time.backup

# deploy new weights
cp weights/gemma2-2b_k-tuned_4_epoch17.lora.h5 ../k-mail-replier/k_mail_replier/weights/gemma2-2b_k-tuned.lora.h5