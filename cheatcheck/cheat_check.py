# -*- coding: utf-8 -*-

import os
import sys
import numpy as np

N = 10
AMP = 28
# DIMEN_TH = [1,3,6,9,12]
VALUE_TH = [float('inf'),3000,2600,2000,1600,800,28]

def dynamic_cluster(data):
	cluster = 1
	for i in range(1,len(data)):
		is_single_cluster = True
		for j in range(i):
			if abs(data[i] - data[j]) <= N:
				is_single_cluster = False
		if is_single_cluster:
			cluster = cluster + 1
	return cluster

def amp_discriminate(data):
	if max(data) - min(data) < AMP:
		return True
	return False

def cheat_discriminate(data):
	if dynamic_cluster(data) == 1:
		if amp_discriminate(data):
			return True
	return False


def value_same_discriminate(data):
	for i in range(1,len(data)):
		if data[i] != data[0]:
			 return False
	return True

#index:1,2,3,4,5,6
def value_range_discriminate(data,index):
	for i in range(len(data)):
		if((data[i] >= VALUE_TH[index - 1]) or (data[i] < VALUE_TH[index])):
			return False
	return True

#n:1,3,6,9,12
def cheat_discriminate_n(data,n):
	if n == 1:
		if value_range_discriminate(data,n):
			return True
	elif n == 3:
		if value_range_discriminate(data,n - 1):
			return True
		elif value_range_discriminate(data,n):
			if not value_same_discriminate(data):
				if cheat_discriminate(data):
					return True
	else:
		is_range = False
		if n == 6:
			if value_range_discriminate(data,4):
				is_range = True
		elif n == 9:
			if value_range_discriminate(data,5):
				is_range = True
		elif n == 12:
			if value_range_discriminate(data,6):
				is_range = True
		if is_range:
			if not value_same_discriminate(data):
				if cheat_discriminate(data):
					return True
	return False

if __name__ == '__main__':
	folder = './test_data'
	for parent,dirnames,filenames in os.walk(folder):
	    for filename in filenames:
	        filedir = os.path.join(parent,filename)
	        data = np.loadtxt(filedir)
	        print("filename:" + filename)
	        print("data:" + str(data))
	        print("all step is:" + str(sum(data)))
	        cheat_index = []
	        for i in range(len(data)):
	        	if i > 10:
	        		data_seg_1 = data[i : i + 1]
	        		data_seg_3 = data[i - 2 : i + 1]
	        		data_seg_6 = data[i - 5 : i + 1]
	        		data_seg_9 = data[i - 8 : i + 1]
	        		data_seg_12 = data[i - 11 : i + 1]
	        		if cheat_discriminate_n(data_seg_12,12):
	        			cheat_index.append([i - 11, i])
	        		elif cheat_discriminate_n(data_seg_9,9):
	        			cheat_index.append([i - 8, i])
	        		elif cheat_discriminate_n(data_seg_6,6):
	        			cheat_index.append([i - 5, i])
	        		elif cheat_discriminate_n(data_seg_3,3):
	        			cheat_index.append([i - 2, i])
	        		elif cheat_discriminate_n(data_seg_1,1):
	        			cheat_index.append([i, i])
	        	elif i > 7:
	        		data_seg_1 = data[i : i + 1]
	        		data_seg_3 = data[i - 2 : i + 1]
	        		data_seg_6 = data[i - 5 : i + 1]
	        		data_seg_9 = data[i - 8 : i + 1]
	        		if cheat_discriminate_n(data_seg_9,9):
	        			cheat_index.append([i - 8, i])
	        		elif cheat_discriminate_n(data_seg_6,6):
	        			cheat_index.append([i - 5, i])
	        		elif cheat_discriminate_n(data_seg_3,3):
	        			cheat_index.append([i - 2, i])
	        		elif cheat_discriminate_n(data_seg_1,1):
	        			cheat_index.append([i, i])
	        	elif i > 4:
	        		data_seg_1 = data[i : i + 1]
	        		data_seg_3 = data[i - 2 : i + 1]
	        		data_seg_6 = data[i - 5 : i + 1]
	        		if cheat_discriminate_n(data_seg_6,6):
	        			cheat_index.append([i - 5, i])
	        		elif cheat_discriminate_n(data_seg_3,3):
	        			cheat_index.append([i - 2, i])
	        		elif cheat_discriminate_n(data_seg_1,1):
	        			cheat_index.append([i, i])
	        	elif i > 1:
	        		data_seg_1 = data[i : i + 1]
	        		data_seg_3 = data[i - 2 : i + 1]
	        		if cheat_discriminate_n(data_seg_3,3):
	        			cheat_index.append([i - 2, i])
	        		elif cheat_discriminate_n(data_seg_1,1):
	        			cheat_index.append([i, i])
	        	else:
	        		data_seg_1 = data[i : i + 1]
	        		if cheat_discriminate_n(data_seg_1,1):
	        			cheat_index.append([i, i])

	        print("cheat data index:" + str(cheat_index))
	        for j in range(len(cheat_index)):
	        	data_seg = data[cheat_index[j][0]:cheat_index[j][1] + 1]
	        	print("cheat data " + str(j) + ":" + str(data_seg))
	        print("\n\n\n")