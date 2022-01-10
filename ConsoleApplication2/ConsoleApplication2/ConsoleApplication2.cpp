// 1dalis.cpp : This file contains the 'main' function. Program execution begins and ends there.
// numeris^2 + numeris^3

#include <iostream>
#include <omp.h>

using namespace std;


class Monitor {
private:
    omp_lock_t* lock;
    int* skaiciukai;
    int dabartinisKiekis;
    int maxKiekis;

public:
    Monitor(omp_lock_t* lock, int size) {
        this->dabartinisKiekis = 0;
        this->skaiciukai = new int[size];
        this->lock = lock;
        this->maxKiekis = size;
    }

    void Add(int skaicius) {
        omp_set_lock(this->lock);
        if (this->dabartinisKiekis >= this->maxKiekis)
        {
            omp_unset_lock(this->lock);
        }

        else {
            this->skaiciukai[this->maxKiekis - this->dabartinisKiekis] = skaicius;
            this->dabartinisKiekis++;
            omp_unset_lock(this->lock);
        }
    }

    void Spausdinti() {
        for (size_t i = 0; i < this->maxKiekis; i++)
        {
            cout << this->skaiciukai[i] << endl;
        }

    }

    int kiek(int kelintas) {
        return skaiciukai[kelintas];
    }

};



int main()
{
    int threduKiekis = 10;
    omp_lock_t RezultataiLock;
    omp_init_lock(&RezultataiLock);
    int sum = 0;
    Monitor Rezultatai = Monitor(&RezultataiLock, threduKiekis);

    omp_set_num_threads(threduKiekis);
#pragma omp parallel
    {
        int numeris = omp_get_thread_num();

        Rezultatai.Add(numeris * numeris + numeris * numeris * numeris);
    }

    omp_destroy_lock(&RezultataiLock);
    Rezultatai.Spausdinti();
#pragma omp parallel for reduction (+:sum)
    for (int i = 1; i < threduKiekis; i++)
        sum = sum + Rezultatai.kiek(i);

    cout << sum;
}
