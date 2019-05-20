package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mpolski/oneview-golang/ov"
)

func main() {

	var (
		clientOV *ov.OVClient

		mEncStatus = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_enclosures_status",
				Help: "Number of enclosures reported given status.",
			},
			[]string{"resourceType", "status"})

		mEncCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_enclosures_count",
				Help: "Total number of enclosures.",
			},
			[]string{"resourceType"})

		mEncAmbTemp = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_enclosures_ambientTemperature_celcius",
				Help: "Ambient temperature as reported by enclosure.",
			},
			[]string{"category", "uuid", "name"})
		mEncAvePwr = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_enclosures_averagePower_watts",
				Help: "Average Power consumption as reported by enclosure.",
			},
			[]string{"category", "uuid", "name"})
		mEncPeakPwr = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_enclosures_peakPower_watts",
				Help: "Peak Power consumption as reported by enclosure.",
			},
			[]string{"category", "uuid", "name"})

		mSrvStatus = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_serverHardware_status",
				Help: "Number of servers reported given status.",
			},
			[]string{"resourceType", "status"})

		mSrvCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_serverHardware_count",
				Help: "Total number of server hardware.",
			},
			[]string{"resourceType"})

		mSrvAmbTemp = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_server_ambientTemperature_celcius",
				Help: "Ambient temperature as reported by server.",
			},
			[]string{"category", "uuid", "name"})
		mSrvAvePwr = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_server_averagePower_watts",
				Help: "Average Power consumption as reported by server.",
			},
			[]string{"category", "uuid", "name"})
		mSrvAveCPU = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_server_CPUFrequency_Mhz",
				Help: "Server CPU frequencyd.",
			},
			[]string{"category", "uuid", "name"})
		mSrvCPUUti = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_server_CPUUtilization_percentage",
				Help: "Server CPU utilization.",
			},
			[]string{"category", "uuid", "name"})
		mSrvPeakPwr = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_server_peakPower_watts",
				Help: "Peak Power consumption as reported by server.",
			},
			[]string{"category", "uuid", "name"})
		mIntCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_interconnect_count",
				Help: "Total number of Interconnects.",
			},
			[]string{"resourceType"})
		mSasIntCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_sasInterconnect_count",
				Help: "Total number of SAS Interconnects.",
			},
			[]string{"resourceType"})

		mSasIntStatus = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "oneview_sasInterconnect_status",
				Help: "Number of SAS Interconnect modules reporting given status.",
			},
			[]string{"resourceType", "status"})
	)

	prometheus.MustRegister(
		mEncCount, mEncStatus, mEncAmbTemp, mEncAvePwr, mEncPeakPwr,
		mSrvCount, mSrvStatus, mSrvAmbTemp, mSrvAvePwr, mSrvAveCPU, mSrvCPUUti, mSrvPeakPwr,
		mIntCount, //mIntStatus,
		mSasIntCount, mSasIntStatus,
	)

	//Enclosure Count, ServerHardware Count, Interconnect Count, SAS Interconnect Count
	go func() {
		ovc := clientOV.NewOVClient(
			os.Getenv("OV_USERNAME"),
			os.Getenv("OV_PASSWORD"),
			os.Getenv("OV_AUTHLOGINDOMAIN"),
			os.Getenv("OV_ENDPOINT"), //use http(s):// prefix
			false,
			800, //why is this a fixed value? either detect or use env var
			"*")
		for {
			//Enclosures count
			encC, err := ovc.GetEnclosures("", "", "", "", "")
			if err != nil {
				fmt.Println("Enclosures Retrieval Failed: ", err)
			} else {
				fmt.Println("#----------------Enclosures Count Total---------------#", encC.Total)
				mEncCount.With(prometheus.Labels{"resourceType": encC.Category}).Set(float64(encC.Total))
			}

			//ServerHardware count
			srvC, err := ovc.GetServerHardwareList([]string{}, "")
			if err != nil {
				fmt.Println("Server Hardware Retrieval Failed: ", err)
			} else {
				fmt.Println("#----------------Server Hardware Count Total---------------#", srvC.Total)
				mSrvCount.With(prometheus.Labels{"resourceType": srvC.Category}).Set(float64(srvC.Total))
			}

			//Interconnect count
			intC, err := ovc.GetInterconnects("", "", "", "")
			if err != nil {
				fmt.Println("Interconnect Retrieval Failed: ", err)
			} else {
				fmt.Println("#----------------Interconnect Count Total---------------#", intC.Total)
				mIntCount.With(prometheus.Labels{"resourceType": intC.Category}).Set(float64(intC.Total))
			}

			//SAS Interconnect count
			sintC, err := ovc.GetSasInterconnects("", []string{}, "", "", "")
			if err != nil {
				fmt.Println("SAS Interconnect Retrieval Failed: ", err)
			} else {
				fmt.Println("#----------------SAS Interconnect Count Total---------------#", sintC.Total)
				mSasIntCount.With(prometheus.Labels{"resourceType": sintC.Category}).Set(float64(sintC.Total))
			}

			fmt.Println(time.Now().Format(time.RFC3339))
			time.Sleep(time.Second * time.Duration(60))
		}
	}()

	//Enclosures status and ServerHardware status
	go func() {
		var (
			f = map[string][]string{
				"OK":       []string{"'status'='OK'"},
				"Unknown":  []string{"'status'='Unknown'"},
				"Warning":  []string{"'status'='Warning'"},
				"Critical": []string{"'status'='Critical'"},
				"Disabled": []string{"'status'='Disabled'"},
			}
		)
		ovc := clientOV.NewOVClient(
			os.Getenv("OV_USERNAME"),
			os.Getenv("OV_PASSWORD"),
			os.Getenv("OV_AUTHLOGINDOMAIN"),
			os.Getenv("OV_ENDPOINT"), //use http(s):// prefix
			false,
			800, //why is this a fixed value? either detect or use env var
			"*")

		for {

			for _, s := range f {

				//Enclsolure status
				encS, err := ovc.GetEnclosures("", "", s[0], "", "")
				if err != nil {
					fmt.Println("Enclosure Retrieval Failed: ", err)
				} else {
					fmt.Println("#----------------Enclosure Count with", s[0], "---------------#", encS.Total)
					mEncStatus.With(prometheus.Labels{"resourceType": encS.Category, "status": s[0]}).Set(float64(encS.Total))
				}
				//ServerHardware status
				srvS, err := ovc.GetServerHardwareList(s, "")
				if err != nil {
					fmt.Println("Server Hardware Retrieval Failed: ", err)
				} else {
					fmt.Println("#----------------Server Hardware Count with", s[0], "---------------#", srvS.Total)
					mSrvStatus.With(prometheus.Labels{"resourceType": srvS.Category, "status": s[0]}).Set(float64(srvS.Total))
				}
				//Interconnects status - not currently supported - needs a workaround
				// intS, err := ovc.GetInterconnects("", "", "'OK'", "")
				// if err != nil {
				// 	fmt.Println("Interconnect Retrieval Failed: ", err)
				// } else {
				// 	fmt.Println("#----------------Interconnect Count with", s[0], "---------------#", intS.Total)
				// 	mSrvStatus.With(prometheus.Labels{"resourceType": intS.Category, "status": "OK"}).Set(float64(intS.Total))
				// }

				//SAS Interconnects status
				sintS, err := ovc.GetSasInterconnects("", s, "", "", "")
				if err != nil {
					fmt.Println("SAS Interconnect Retrieval Failed: ", err)
				} else {
					fmt.Println("#----------------SAS Interconnect Count with", s[0], "---------------#", sintS.Total)
					mSasIntStatus.With(prometheus.Labels{"resourceType": sintS.Category, "status": s[0]}).Set(float64(sintS.Total))
				}
			}

			fmt.Println(time.Now().Format(time.RFC3339))
			time.Sleep(time.Second * time.Duration(60))
		}

	}()

	// //GetUtilization for Enclosure
	go func() {
		ovc := clientOV.NewOVClient(
			os.Getenv("OV_USERNAME"),
			os.Getenv("OV_PASSWORD"),
			os.Getenv("OV_AUTHLOGINDOMAIN"),
			os.Getenv("OV_ENDPOINT"), //use http(s):// prefix
			false,
			800, //why is this a fixed value? either detect or use env var
			"*")

		for {
			c, err := ovc.GetEnclosures("", "", "", "", "")
			if err != nil {
				fmt.Println("Enclosures List Retrieval Failed: ", err)
			} else {
				for i := 0; i < len(c.Members); i++ {
					category := c.Category
					UUID := c.Members[i].UUID
					name := c.Members[i].Name

					w, err := ovc.GetUtilization("", "", "true", "", c.Members[i].URI.String()) // "true" should be boolean, currently SetQueryString (/rest/netutil.go) accepts map of strings?
					if err != nil {
						fmt.Println("Enclosure utilization Data Retrieval Failed: ,", err)
					}
					//if data not fresh, set values to 0
					if w.IsFresh == false {
						fmt.Println("No fresh utilization data vailable for -", name, "-...skipping")
						mEncAmbTemp.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
						mEncAvePwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
						mEncPeakPwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
					} else {
						for j := 0; j < len(w.MetricList); j++ {
							mEncAmbTemp.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(w.MetricList[0].MetricSamples[0][1].(float64))
							mEncAvePwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(w.MetricList[1].MetricSamples[0][1].(float64))
							mEncPeakPwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(w.MetricList[2].MetricSamples[0][1].(float64))

							fmt.Println("Category:", category, "UUID: ", UUID, ", Name: ", name, ", Metric:", w.MetricList[j].MetricName, ", Value: ", w.MetricList[j].MetricSamples[0][1], ", Sample Time: ", w.MetricList[j].MetricSamples[0][0])
						}
					}
				}
			}
			fmt.Println(time.Now().Format(time.RFC3339))
			time.Sleep(time.Second * time.Duration(300))
		}
	}()

	//GetUtilization for ServerHardware
	go func() {
		ovc := clientOV.NewOVClient(
			os.Getenv("OV_USERNAME"),
			os.Getenv("OV_PASSWORD"),
			os.Getenv("OV_AUTHLOGINDOMAIN"),
			os.Getenv("OV_ENDPOINT"), //use http(s):// prefix
			false,
			800, //why is this a fixed value? either detect or use env var
			"*")
		for {
			filters := []string{}
			c, err := ovc.GetServerHardwareList(filters, "")
			if err != nil {
				fmt.Println("Server Hardware List Retrieval Failed: ", err)
			}

			for i := 0; i < len(c.Members); i++ {
				category := c.Category
				UUID := c.Members[i].UUID
				name := c.Members[i].Name

				u, err := ovc.GetUtilization("", "", "true", "", c.Members[i].URI.String()) // "true" should be boolean, currently SetQueryString (/rest/netutil.go) accepts map of strings?
				if err != nil {
					fmt.Println("Server Hardware utilization Data Retrieval Failed: ,", err)
				}
				//if data not fresh, set values to 0
				if u.IsFresh == false {
					fmt.Println("No fresh utilization data available for -", name, "-...skipping")
					mSrvAmbTemp.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
					mSrvAvePwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
					mSrvAveCPU.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
					mSrvCPUUti.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
					mSrvPeakPwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(0)
				} else {
					for j := 0; j < len(u.MetricList); j++ {
						mSrvAmbTemp.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(u.MetricList[0].MetricSamples[0][1].(float64))
						mSrvAvePwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(u.MetricList[1].MetricSamples[0][1].(float64))
						mSrvAveCPU.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(u.MetricList[2].MetricSamples[0][1].(float64))
						mSrvCPUUti.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(u.MetricList[3].MetricSamples[0][1].(float64))
						mSrvPeakPwr.With(prometheus.Labels{"category": category, "uuid": UUID, "name": name}).Set(u.MetricList[4].MetricSamples[0][1].(float64))

						fmt.Println("Category:", category, "UUID: ", UUID, ", Name: ", name, ", Metric:", u.MetricList[j].MetricName, ", Value: ", u.MetricList[j].MetricSamples[0][1], ", Sample Time: ", u.MetricList[j].MetricSamples[0][0])
					}
				}
			}
			fmt.Println(time.Now().Format(time.RFC3339))
			time.Sleep(time.Second * time.Duration(300))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	//log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	http.ListenAndServe(":8080", nil)
}
