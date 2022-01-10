#include <iostream>
#include <omp.h>
using namespace std;

int main() {
	// pradiniai duomenys
	int* arr = new int[1000];
	float* rez = new float[10];
	int* sumos = new int[10];
	int sum = 0;
	int kiek = -1;
	for (int i = 0; i < 999; i++) {
		arr[i] = i + 1;
		if (i < 10) {
			rez[i] = -1;
		}
	}
	// giju nustatymas
	omp_set_num_threads(10);
	int threadNumber = 0;
#pragma omp parallel private(threadNumber)
	{
		threadNumber = omp_get_thread_num();
		int startingIndex = threadNumber * 100; 	// giju padalinimas
		int endIndex = startingIndex + 99;
		float vidurkis = 0;
		int daliklis = 99;
		int suma = 0;
		for (int i = startingIndex; i < endIndex; i++) {
			suma += arr[i];
		}
		vidurkis = suma / daliklis;

		kiek++;
		rez[kiek] = vidurkis;
		sumos[kiek]=suma;

		//cout << vidurkis << endl;
	}
#pragma omp parallel for reduction (+:sum)
	for (int i = 0; i < 10; i++)
		sum = sum + sumos[i];

	cout << sum;
	for (int i = 0; i < 10; i++)
	{
		cout << rez[i] << endl;
	}
}